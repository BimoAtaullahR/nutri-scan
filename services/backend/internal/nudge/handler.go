package nudge

import (
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const (
	ActionEatAsPlanned    = "eat_as_planned"
	ActionSetAsidePortion = "set_aside_portion"
	ActionReviewFood      = "review_food"

	ConfidenceLevelLow    = "low"
	ConfidenceLevelMedium = "medium"
	ConfidenceLevelHigh   = "high"

	BasisGeneric      = "generic"
	BasisPersonalized = "personalized"
	BasisReviewFood   = "review_food"
)

type Decision struct {
	NudgeID                      string   `json:"nudgeId"`
	Action                       string   `json:"action"`
	Message                      string   `json:"message"`
	EstimatedPreventedEnergyKcal *int     `json:"estimatedPreventedEnergyKcal"`
	ConfidenceLevel              string   `json:"confidenceLevel"`
	IsPersonalized               bool     `json:"isPersonalized"`
	Basis                        string   `json:"basis"`
	DisplayTags                  []string `json:"displayTags"`
}

type DecisionInput struct {
	FoodCategory        string
	FoodConfidence      float64
	IsLowConfidence     bool
	EstimatedEnergyKcal int
	HasUserProfile      bool
	BMICategory         string
}

func Decide(input DecisionInput) Decision {
	decision := Decision{
		NudgeID:         uuid.NewString(),
		ConfidenceLevel: confidenceLevel(input.FoodConfidence),
		DisplayTags:     []string{},
	}

	if input.IsLowConfidence {
		decision.Action = ActionReviewFood
		decision.Message = "Review this food result before using it for your plan."
		decision.Basis = BasisReviewFood
		decision.DisplayTags = []string{"Review food", "Low confidence"}
		return decision
	}

	decision.Action = ActionEatAsPlanned
	decision.Message = "This portion looks reasonable to eat as planned."
	decision.Basis = BasisGeneric
	decision.DisplayTags = []string{"Estimated energy", displayFoodCategory(input.FoodCategory)}

	if input.EstimatedEnergyKcal >= 600 {
		decision.Action = ActionSetAsidePortion
		decision.Message = "Consider setting aside a smaller portion before eating."
		preventedEnergy := int(math.Round(float64(input.EstimatedEnergyKcal) * 0.25))
		decision.EstimatedPreventedEnergyKcal = &preventedEnergy
		decision.DisplayTags = append(decision.DisplayTags, "Portion nudge")
	}

	if input.HasUserProfile {
		decision.IsPersonalized = true
		decision.Basis = BasisPersonalized
		decision.DisplayTags = append(decision.DisplayTags, "Personalized")
		switch input.BMICategory {
		case "overweight", "obese":
			if decision.Action == ActionEatAsPlanned && input.EstimatedEnergyKcal >= 450 {
				decision.Action = ActionSetAsidePortion
				decision.Message = "Based on your profile, consider setting aside a small portion before eating."
				preventedEnergy := int(math.Round(float64(input.EstimatedEnergyKcal) * 0.20))
				decision.EstimatedPreventedEnergyKcal = &preventedEnergy
			}
		case "underweight":
			if decision.Action == ActionSetAsidePortion && input.EstimatedEnergyKcal < 750 {
				decision.Action = ActionEatAsPlanned
				decision.Message = "Based on your profile, this portion looks reasonable to eat as planned."
				decision.EstimatedPreventedEnergyKcal = nil
			}
		}
	}

	return decision
}

func confidenceLevel(score float64) string {
	switch {
	case score >= 0.8:
		return ConfidenceLevelHigh
	case score >= 0.6:
		return ConfidenceLevelMedium
	default:
		return ConfidenceLevelLow
	}
}

func displayFoodCategory(foodCategory string) string {
	if foodCategory == "" {
		return "Food estimate"
	}

	return foodCategory
}

type Handler struct {
	logger *slog.Logger
	store  Store
}

func NewHandler(store Store, logger *slog.Logger) *Handler {
	return &Handler{
		logger: logger,
		store:  store,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, requireAnonymousUser func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(requireAnonymousUser)
		r.Post("/nudges/{nudgeID}/responses", h.recordResponse)
	})
}

func (h *Handler) recordResponse(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := user.AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	nudgeID := chi.URLParam(r, "nudgeID")
	if _, err := uuid.Parse(nudgeID); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "nudge_not_found",
			"message": "nudge was not found",
		})
		return
	}

	var request recordNudgeResponseRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_nudge_response_request",
			"message": "nudge response request must be valid JSON",
		})
		return
	}

	if validationErr := request.validate(); validationErr != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_nudge_response_request",
			"message": validationErr,
		})
		return
	}

	err := h.store.VerifyNudgeOwnership(r.Context(), anonymousUser.ID, request.ScanID, nudgeID)
	if errors.Is(err, ErrNudgeNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "nudge_not_found",
			"message": "nudge was not found",
		})
		return
	}
	if err != nil {
		h.logger.Error("failed to verify nudge ownership", "nudge_id", nudgeID, "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "nudge_verification_failed",
			"message": "nudge ownership could not be verified",
		})
		return
	}

	record, err := h.store.RecordResponse(r.Context(), ResponseRecord{
		ID:              uuid.NewString(),
		ScanID:          request.ScanID,
		AnonymousUserID: anonymousUser.ID,
		Response:        request.Response,
	})
	if err != nil {
		h.logger.Error("failed to record nudge response", "nudge_id", nudgeID, "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "nudge_response_record_failed",
			"message": "nudge response could not be recorded",
		})
		return
	}

	writeJSON(w, http.StatusCreated, newNudgeResponseResponse(record, nudgeID))
}

type recordNudgeResponseRequest struct {
	ScanID   string `json:"scanId"`
	Response string `json:"response"`
}

func (r recordNudgeResponseRequest) validate() string {
	if _, err := uuid.Parse(r.ScanID); err != nil {
		return "scanId must be a valid UUID"
	}

	switch r.Response {
	case "followed", "did_not_follow", "dismissed":
		return ""
	default:
		return "response must be one of followed, did_not_follow, or dismissed"
	}
}

type nudgeResponseResponse struct {
	NudgeResponseID string    `json:"nudgeResponseId"`
	NudgeID         string    `json:"nudgeId"`
	ScanID          string    `json:"scanId"`
	Response        string    `json:"response"`
	CreatedAt       time.Time `json:"createdAt"`
}

func newNudgeResponseResponse(record ResponseRecord, nudgeID string) nudgeResponseResponse {
	return nudgeResponseResponse{
		NudgeResponseID: record.ID,
		NudgeID:         nudgeID,
		ScanID:          record.ScanID,
		Response:        record.Response,
		CreatedAt:       record.CreatedAt,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
