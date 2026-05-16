package summary

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
)

const defaultDailyGoalEnergyKcal = 2000

type Handler struct {
	store  Store
	logger *slog.Logger
	clock  func() time.Time
}

func NewHandler(store Store, logger *slog.Logger) *Handler {
	return &Handler{
		store:  store,
		logger: logger,
		clock:  time.Now,
	}
}

func (h *Handler) RegisterRoutes(r chi.Router, requireAnonymousUser func(http.Handler) http.Handler) {
	r.Group(func(r chi.Router) {
		r.Use(requireAnonymousUser)
		r.Get("/summaries/daily", h.getDailySummary)
	})
}

func (h *Handler) getDailySummary(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := user.AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	dateStr := r.URL.Query().Get("date")
	var targetDate time.Time
	var err error

	if dateStr == "" {
		targetDate = h.clock().Local().Truncate(24 * time.Hour)
		dateStr = targetDate.Format("2006-01-02")
	} else {
		targetDate, err = time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error":   "invalid_date_format",
				"message": "date must be in YYYY-MM-DD format",
			})
			return
		}
	}

	dateStart := targetDate
	dateEnd := targetDate.Add(24 * time.Hour)

	mealSummaries, err := h.store.GetDailyMealSummaries(r.Context(), anonymousUser.ID, dateStart, dateEnd)
	if err != nil {
		h.logger.Error("failed to get daily meal summaries", "anonymous_user_id", anonymousUser.ID, "date", dateStr, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "summary_fetch_failed",
			"message": "failed to fetch daily summary",
		})
		return
	}

	summary := buildDailyEnergySummary(dateStr, mealSummaries)
	writeJSON(w, http.StatusOK, summary)
}

type mealEnergySummaryResponse struct {
	MealType        string `json:"mealType"`
	EatenEnergyKcal int    `json:"eatenEnergyKcal"`
	ScanCount       int    `json:"scanCount"`
}

type dailyEnergySummaryResponse struct {
	Date                string                      `json:"date"`
	EatenEnergyKcal     int                         `json:"eatenEnergyKcal"`
	RemainingEnergyKcal int                         `json:"remainingEnergyKcal"`
	BurnedEnergyKcal    int                         `json:"burnedEnergyKcal"`
	DailyGoalEnergyKcal int                         `json:"dailyGoalEnergyKcal"`
	Meals               []mealEnergySummaryResponse `json:"meals"`
}

func buildDailyEnergySummary(dateStr string, dbSummaries []MealSummary) dailyEnergySummaryResponse {
	mealsMap := map[string]MealSummary{
		"breakfast": {MealType: "breakfast", EatenEnergyKcal: 0, ScanCount: 0},
		"lunch":     {MealType: "lunch", EatenEnergyKcal: 0, ScanCount: 0},
		"dinner":    {MealType: "dinner", EatenEnergyKcal: 0, ScanCount: 0},
		"snack":     {MealType: "snack", EatenEnergyKcal: 0, ScanCount: 0},
	}

	for _, s := range dbSummaries {
		mealsMap[s.MealType] = s
	}

	orderedMeals := []mealEnergySummaryResponse{
		{MealType: "breakfast", EatenEnergyKcal: mealsMap["breakfast"].EatenEnergyKcal, ScanCount: mealsMap["breakfast"].ScanCount},
		{MealType: "lunch", EatenEnergyKcal: mealsMap["lunch"].EatenEnergyKcal, ScanCount: mealsMap["lunch"].ScanCount},
		{MealType: "dinner", EatenEnergyKcal: mealsMap["dinner"].EatenEnergyKcal, ScanCount: mealsMap["dinner"].ScanCount},
		{MealType: "snack", EatenEnergyKcal: mealsMap["snack"].EatenEnergyKcal, ScanCount: mealsMap["snack"].ScanCount},
	}

	totalEaten := 0
	for _, m := range orderedMeals {
		totalEaten += m.EatenEnergyKcal
	}

	goal := defaultDailyGoalEnergyKcal
	burned := 0
	remaining := goal - totalEaten + burned

	return dailyEnergySummaryResponse{
		Date:                dateStr,
		EatenEnergyKcal:     totalEaten,
		RemainingEnergyKcal: remaining,
		BurnedEnergyKcal:    burned,
		DailyGoalEnergyKcal: goal,
		Meals:               orderedMeals,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
