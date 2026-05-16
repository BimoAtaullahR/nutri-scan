package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"math"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

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

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/anonymous-users", h.createAnonymousUser)
	r.Group(func(r chi.Router) {
		r.Use(h.RequireAnonymousUser)
		r.Get("/me/profile", h.getProfile)
		r.Put("/me/profile", h.updateProfile)
	})
}

func (h *Handler) createAnonymousUser(w http.ResponseWriter, r *http.Request) {
	tokenSuffix, err := randomHex(32)
	if err != nil {
		h.logger.Error("failed to generate anonymous user token", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_create_failed",
			"message": "anonymous user could not be created",
		})
		return
	}

	token := "nutriscan_" + tokenSuffix
	anonymousUser := AnonymousUser{
		ID:        uuid.NewString(),
		TokenHash: hashBearerToken(token),
	}

	if err := h.store.CreateAnonymousUser(r.Context(), anonymousUser); err != nil {
		h.logger.Error("failed to create anonymous user", "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_create_failed",
			"message": "anonymous user could not be created",
		})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"anonymousUserId": anonymousUser.ID,
		"accessToken":     token,
		"tokenType":       "Bearer",
	})
}

func (h *Handler) RequireAnonymousUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := bearerToken(r.Header.Get("Authorization"))
		if !ok {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error":   "missing_bearer_token",
				"message": "anonymous user bearer token is required",
			})
			return
		}

		anonymousUser, err := h.store.GetAnonymousUserByTokenHash(r.Context(), hashBearerToken(token))
		if errors.Is(err, ErrAnonymousUserNotFound) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error":   "invalid_bearer_token",
				"message": "anonymous user bearer token is invalid",
			})
			return
		}
		if err != nil {
			h.logger.Error("failed to authenticate anonymous user", "error", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error":   "anonymous_user_auth_failed",
				"message": "anonymous user could not be authenticated",
			})
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), anonymousUserContextKey{}, anonymousUser)))
	})
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	profile, err := h.store.GetUserProfile(r.Context(), anonymousUser.ID)
	if errors.Is(err, ErrUserProfileNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error":   "user_profile_not_found",
			"message": "user profile has not been created",
		})
		return
	}
	if err != nil {
		h.logger.Error("failed to get user profile", "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "user_profile_get_failed",
			"message": "user profile could not be loaded",
		})
		return
	}

	writeJSON(w, http.StatusOK, newUserProfileResponse(profile))
}

func (h *Handler) updateProfile(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	var request updateUserProfileRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_user_profile_request",
			"message": "user profile request must be valid JSON",
		})
		return
	}

	if validationErr := request.validate(); validationErr != "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid_user_profile_request",
			"message": validationErr,
		})
		return
	}

	bmi := calculateBMI(request.HeightCm, request.WeightKg)
	profile, err := h.store.UpsertUserProfile(r.Context(), UserProfile{
		AnonymousUserID: anonymousUser.ID,
		HeightCm:        request.HeightCm,
		WeightKg:        request.WeightKg,
		AgeRange:        request.ageRange(),
		BMI:             bmi,
		BMICategory:     bmiCategory(bmi),
	})
	if err != nil {
		h.logger.Error("failed to upsert user profile", "anonymous_user_id", anonymousUser.ID, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "user_profile_update_failed",
			"message": "user profile could not be saved",
		})
		return
	}

	writeJSON(w, http.StatusOK, newUserProfileResponse(profile))
}

type updateUserProfileRequest struct {
	HeightCm float64 `json:"heightCm"`
	WeightKg float64 `json:"weightKg"`
	AgeRange *string `json:"ageRange"`
}

func (r updateUserProfileRequest) validate() string {
	if r.HeightCm <= 0 {
		return "heightCm must be greater than 0"
	}
	if r.HeightCm > 300 {
		return "heightCm must be less than or equal to 300"
	}
	if r.WeightKg <= 0 {
		return "weightKg must be greater than 0"
	}
	if r.WeightKg > 500 {
		return "weightKg must be less than or equal to 500"
	}
	if r.AgeRange != nil && !validAgeRange(*r.AgeRange) {
		return "ageRange must be one of under_18, 18_24, 25_34, 35_44, 45_54, or 55_plus"
	}

	return ""
}

func (r updateUserProfileRequest) ageRange() string {
	if r.AgeRange == nil {
		return ""
	}

	return *r.AgeRange
}

func validAgeRange(ageRange string) bool {
	switch ageRange {
	case "under_18", "18_24", "25_34", "35_44", "45_54", "55_plus":
		return true
	default:
		return false
	}
}

type userProfileResponse struct {
	HeightCm    float64 `json:"heightCm"`
	WeightKg    float64 `json:"weightKg"`
	AgeRange    string  `json:"ageRange,omitempty"`
	BMI         float64 `json:"bmi"`
	BMICategory string  `json:"bmiCategory"`
}

func newUserProfileResponse(profile UserProfile) userProfileResponse {
	return userProfileResponse{
		HeightCm:    profile.HeightCm,
		WeightKg:    profile.WeightKg,
		AgeRange:    profile.AgeRange,
		BMI:         profile.BMI,
		BMICategory: profile.BMICategory,
	}
}

func calculateBMI(heightCm float64, weightKg float64) float64 {
	heightM := heightCm / 100
	bmi := weightKg / (heightM * heightM)
	return math.Round(bmi*100) / 100
}

func bmiCategory(bmi float64) string {
	switch {
	case bmi < 18.5:
		return "underweight"
	case bmi < 25:
		return "normal"
	case bmi < 30:
		return "overweight"
	default:
		return "obese"
	}
}

func randomHex(size int) (string, error) {
	buf := make([]byte, size)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}

type anonymousUserContextKey struct{}

func AnonymousUserFromContext(ctx context.Context) (AnonymousUser, bool) {
	anonymousUser, ok := ctx.Value(anonymousUserContextKey{}).(AnonymousUser)
	return anonymousUser, ok
}

func bearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") || token == "" {
		return "", false
	}

	return token, true
}

func hashBearerToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func writePlaceholder(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": message,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
