package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"math"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/nudge"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const maxScanImageBytes = 8 << 20

type Handler struct {
	aiClient         InferenceClient
	clock            func() time.Time
	inferenceTimeout time.Duration
	logger           *slog.Logger
	profileReader    ProfileReader
	store            Store
}

type ProfileReader interface {
	GetUserProfile(ctx context.Context, anonymousUserID string) (user.UserProfile, error)
}

func NewHandler(store Store, profileReader ProfileReader, aiClient InferenceClient, logger *slog.Logger) *Handler {
	return &Handler{
		aiClient:         aiClient,
		clock:            time.Now,
		inferenceTimeout: 10 * time.Second,
		logger:           logger,
		profileReader:    profileReader,
		store:            store,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, requireAnonymousUser func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(requireAnonymousUser)
		r.Post("/scans", h.createScan)
		r.Get("/scans/{scanID}", h.getScan)
	})
}

func (h *Handler) createScan(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := user.AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	imageBytes, contentType, err := scanImage(w, r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_scan_image",
			"message": err.Error(),
		})
		return
	}

	mealType := mealTypeFromRequest(r)
	if mealType == "" {
		mealType = fallbackMealType(h.clock())
	}
	if !validMealType(mealType) {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_scan_request",
			"message": "mealType must be one of breakfast, lunch, dinner, or snack",
		})
		return
	}

	scan := Scan{
		ID:              uuid.NewString(),
		AnonymousUserID: anonymousUser.ID,
		Status:          ScanStatusProcessing,
		MealType:        mealType,
	}
	if err := h.store.CreateProcessingScan(r.Context(), scan); err != nil {
		h.logger.Error("failed to create processing scan", "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "scan_create_failed",
			"message": "scan could not be created",
		})
		return
	}

	inferenceCtx, cancel := context.WithTimeout(r.Context(), h.inferenceTimeout)
	defer cancel()

	result, err := h.aiClient.InferScan(inferenceCtx, ScanImage{
		ContentType: contentType,
		Bytes:       imageBytes,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(inferenceCtx.Err(), context.DeadlineExceeded) {
			processingScan, getErr := h.store.GetScan(r.Context(), anonymousUser.ID, scan.ID)
			if getErr != nil {
				h.logger.Error("failed to load processing scan after inference timeout", "scan_id", scan.ID, "error", getErr)
				writeJSON(w, http.StatusInternalServerError, map[string]string{
					"error":   "scan_get_failed",
					"message": "scan could not be loaded",
				})
				return
			}
			writeJSON(w, http.StatusAccepted, newScanResponse(processingScan))
			return
		}

		failedScan, failErr := h.store.FailScan(r.Context(), anonymousUser.ID, scan.ID, "ai_inference_failed")
		if failErr != nil {
			h.logger.Error("failed to mark scan failed", "scan_id", scan.ID, "error", failErr)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "scan_update_failed",
				"message": "scan could not be updated",
			})
			return
		}
		writeJSON(w, http.StatusCreated, newScanResponse(failedScan))
		return
	}

	if validationErr := validateInferenceResult(result); validationErr != "" {
		failedScan, failErr := h.store.FailScan(r.Context(), anonymousUser.ID, scan.ID, validationErr)
		if failErr != nil {
			h.logger.Error("failed to mark scan failed after invalid inference payload", "scan_id", scan.ID, "error", failErr)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "scan_update_failed",
				"message": "scan could not be updated",
			})
			return
		}
		writeJSON(w, http.StatusCreated, newScanResponse(failedScan))
		return
	}

	nudgeDecision, err := h.decideNudge(r.Context(), anonymousUser.ID, result)
	if err != nil {
		h.logger.Error("failed to produce Nudge Decision", "scan_id", scan.ID, "anonymous_user_id", anonymousUser.ID, "error", err)
		failedScan, failErr := h.store.FailScan(r.Context(), anonymousUser.ID, scan.ID, "nudge_decision_failed")
		if failErr != nil {
			h.logger.Error("failed to mark scan failed after Nudge Decision failure", "scan_id", scan.ID, "error", failErr)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "scan_update_failed",
				"message": "scan could not be updated",
			})
			return
		}
		writeJSON(w, http.StatusCreated, newScanResponse(failedScan))
		return
	}

	completedScan, err := h.store.CompleteScan(r.Context(), anonymousUser.ID, scan.ID, result, nudgeDecision)
	if err != nil {
		h.logger.Error("failed to complete scan", "scan_id", scan.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "scan_update_failed",
			"message": "scan could not be updated",
		})
		return
	}

	writeJSON(w, http.StatusCreated, newScanResponse(completedScan))
}

func (h *Handler) decideNudge(ctx context.Context, anonymousUserID string, result InferenceResult) (nudge.Decision, error) {
	var profile user.UserProfile
	hasProfile := false
	if h.profileReader != nil {
		var err error
		profile, err = h.profileReader.GetUserProfile(ctx, anonymousUserID)
		if err == nil {
			hasProfile = true
		} else if !errors.Is(err, user.ErrUserProfileNotFound) {
			return nudge.Decision{}, err
		}
	}

	estimatedEnergyKcal := 0
	if result.EstimatedEnergyRange != nil {
		estimatedEnergyKcal = int(math.Round(float64(result.EstimatedEnergyRange.MinKcal+result.EstimatedEnergyRange.MaxKcal) / 2))
	}

	return nudge.Decide(nudge.DecisionInput{
		FoodCategory:        result.FoodCategory.Slug,
		FoodConfidence:      result.FoodCategory.ConfidenceScore,
		IsLowConfidence:     result.IsLowConfidence,
		EstimatedEnergyKcal: estimatedEnergyKcal,
		HasUserProfile:      hasProfile,
		BMICategory:         profile.BMICategory,
	}), nil
}

func (h *Handler) getScan(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := user.AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	scanID := chi.URLParam(r, "scanID")
	if _, err := uuid.Parse(scanID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "scan_not_found",
			"message": "scan was not found",
		})
		return
	}

	scan, err := h.store.GetScan(r.Context(), anonymousUser.ID, scanID)
	if errors.Is(err, ErrScanNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "scan_not_found",
			"message": "scan was not found",
		})
		return
	}
	if err != nil {
		h.logger.Error("failed to get scan", "scan_id", scanID, "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "scan_get_failed",
			"message": "scan could not be loaded",
		})
		return
	}

	writeJSON(w, http.StatusOK, newScanResponse(scan))
}

func scanImage(w http.ResponseWriter, r *http.Request) ([]byte, string, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxScanImageBytes+1)
	if err := r.ParseMultipartForm(maxScanImageBytes); err != nil {
		return nil, "", errors.New("scan image must be multipart form data with an image field no larger than 8 MB")
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, "", errors.New("image field is required")
	}
	defer file.Close()

	imageBytes, err := readLimited(file, maxScanImageBytes)
	if err != nil {
		return nil, "", err
	}
	if len(imageBytes) == 0 {
		return nil, "", errors.New("image must not be empty")
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(imageBytes)
	}
	if !validScanImageContentType(contentType) {
		return nil, "", errors.New("image must be JPEG, PNG, or WebP")
	}

	return imageBytes, contentType, nil
}

func readLimited(file multipart.File, limit int64) ([]byte, error) {
	var buf bytes.Buffer
	written, err := io.CopyN(&buf, file, limit+1)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.New("image could not be read")
	}
	if written > limit {
		return nil, errors.New("image must be no larger than 8 MB")
	}

	return buf.Bytes(), nil
}

func validScanImageContentType(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	switch contentType {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func mealTypeFromRequest(r *http.Request) string {
	if value := strings.TrimSpace(r.FormValue("mealType")); value != "" {
		return value
	}

	return ""
}

func fallbackMealType(now time.Time) string {
	hour := now.Local().Hour()
	switch {
	case hour >= 5 && hour <= 10:
		return MealTypeBreakfast
	case hour >= 11 && hour <= 15:
		return MealTypeLunch
	case hour >= 16 && hour <= 20:
		return MealTypeDinner
	default:
		return MealTypeSnack
	}
}

func validMealType(mealType string) bool {
	switch mealType {
	case MealTypeBreakfast, MealTypeLunch, MealTypeDinner, MealTypeSnack:
		return true
	default:
		return false
	}
}

func validateInferenceResult(result InferenceResult) string {
	if result.ModelVersion == "" {
		return "invalid_inference_payload"
	}
	if result.FoodCategory.Slug == "" {
		return "invalid_inference_payload"
	}
	if result.FoodCategory.ConfidenceScore < 0 || result.FoodCategory.ConfidenceScore > 1 {
		return "invalid_inference_payload"
	}
	if result.CoarsePortion != "small" && result.CoarsePortion != "medium" && result.CoarsePortion != "large" {
		return "invalid_inference_payload"
	}
	if result.ConfidenceThreshold < 0 || result.ConfidenceThreshold > 1 {
		return "invalid_inference_payload"
	}
	if result.EstimatedEnergyRange != nil {
		if result.EstimatedEnergyRange.MinKcal < 0 || result.EstimatedEnergyRange.MaxKcal < 0 {
			return "invalid_inference_payload"
		}
		if result.EstimatedEnergyRange.MinKcal > result.EstimatedEnergyRange.MaxKcal {
			return "invalid_inference_payload"
		}
	}

	for _, alternative := range result.Alternatives {
		if alternative.Slug == "" || alternative.ConfidenceScore < 0 || alternative.ConfidenceScore > 1 {
			return "invalid_inference_payload"
		}
	}

	return ""
}

type scanResponse struct {
	ScanID               string                `json:"scanId"`
	Status               string                `json:"status"`
	MealType             string                `json:"mealType"`
	EstimatedEnergyKcal  *int                  `json:"estimatedEnergyKcal"`
	EstimatedEnergyRange *estimatedEnergyRange `json:"estimatedEnergyRange"`
	Inference            *inferenceSummary     `json:"inference"`
	NudgeDecision        *nudge.Decision       `json:"nudgeDecision"`
	FailureReason        *string               `json:"failureReason"`
	CreatedAt            time.Time             `json:"createdAt"`
	UpdatedAt            time.Time             `json:"updatedAt"`
}

type estimatedEnergyRange struct {
	MinKcal int `json:"minKcal"`
	MaxKcal int `json:"maxKcal"`
}

type inferenceSummary struct {
	FoodCategory    string                   `json:"foodCategory"`
	FoodConfidence  float64                  `json:"foodConfidence"`
	Alternatives    []foodCategoryConfidence `json:"alternatives"`
	CoarsePortion   string                   `json:"coarsePortion"`
	IsLowConfidence bool                     `json:"isLowConfidence"`
}

type foodCategoryConfidence struct {
	FoodCategory string  `json:"foodCategory"`
	Confidence   float64 `json:"confidence"`
}

func newScanResponse(scan Scan) scanResponse {
	response := scanResponse{
		ScanID:        scan.ID,
		Status:        scan.Status,
		MealType:      scan.MealType,
		FailureReason: stringPtrOrNil(scan.FailureReason),
		CreatedAt:     scan.CreatedAt,
		UpdatedAt:     scan.UpdatedAt,
	}

	if scan.EstimatedEnergyRange != nil {
		response.EstimatedEnergyRange = &estimatedEnergyRange{
			MinKcal: scan.EstimatedEnergyRange.MinKcal,
			MaxKcal: scan.EstimatedEnergyRange.MaxKcal,
		}
		estimatedEnergyKcal := int(math.Round(float64(scan.EstimatedEnergyRange.MinKcal+scan.EstimatedEnergyRange.MaxKcal) / 2))
		response.EstimatedEnergyKcal = &estimatedEnergyKcal
	}

	if scan.Inference != nil {
		alternatives := make([]foodCategoryConfidence, 0, len(scan.Inference.Alternatives))
		for _, alternative := range scan.Inference.Alternatives {
			alternatives = append(alternatives, foodCategoryConfidence{
				FoodCategory: alternative.Slug,
				Confidence:   alternative.ConfidenceScore,
			})
		}
		response.Inference = &inferenceSummary{
			FoodCategory:    scan.Inference.FoodCategory.Slug,
			FoodConfidence:  scan.Inference.FoodCategory.ConfidenceScore,
			Alternatives:    alternatives,
			CoarsePortion:   scan.Inference.CoarsePortion,
			IsLowConfidence: scan.Inference.IsLowConfidence,
		}
	}
	if scan.NudgeDecision != nil {
		response.NudgeDecision = scan.NudgeDecision
	}

	return response
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
