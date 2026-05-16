package nudge

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"

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
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Post("/nudges/{nudgeID}/responses", h.recordResponse)
}

func (h *Handler) recordResponse(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusNotImplemented, map[string]string{
		"error":   "not_implemented",
		"message": "nudge response recording is not implemented yet",
		"nudgeId": chi.URLParam(r, "nudgeID"),
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
