package trend

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
)

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
		r.Get("/trends/weekly", h.getWeeklyTrend)
	})
}

func (h *Handler) getWeeklyTrend(w http.ResponseWriter, r *http.Request) {
	anonymousUser, ok := user.AnonymousUserFromContext(r.Context())
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "anonymous_user_context_missing",
			"message": "anonymous user context is missing",
		})
		return
	}

	dateStr := r.URL.Query().Get("weekStart")
	var targetWeekStart time.Time
	var err error

	if dateStr == "" {
		targetWeekStart = getWeekStart(h.clock().Local())
		dateStr = targetWeekStart.Format("2006-01-02")
	} else {
		targetWeekStart, err = time.ParseInLocation("2006-01-02", dateStr, time.Local)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error":   "invalid_date_format",
				"message": "weekStart must be in YYYY-MM-DD format",
			})
			return
		}
		// ensure it's actually the start of a week (e.g. Monday)
		targetWeekStart = getWeekStart(targetWeekStart)
		dateStr = targetWeekStart.Format("2006-01-02")
	}

	weekEnd := targetWeekStart.AddDate(0, 0, 7)

	dbSummaries, err := h.store.GetWeeklyTrendSummaries(r.Context(), anonymousUser.ID, targetWeekStart, weekEnd)
	if err != nil {
		h.logger.Error("failed to get weekly trend summaries", "anonymous_user_id", anonymousUser.ID, "weekStart", dateStr, "error", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "trend_fetch_failed",
			"message": "failed to fetch weekly trend",
		})
		return
	}

	trend := buildWeeklyEnergyTrend(targetWeekStart, dbSummaries)
	writeJSON(w, http.StatusOK, trend)
}

func getWeekStart(t time.Time) time.Time {
	offset := int(time.Monday - t.Weekday())
	if offset > 0 {
		offset = -6
	}
	return t.AddDate(0, 0, offset).Truncate(24 * time.Hour)
}

type dailyEnergyTrendPointResponse struct {
	Date            string `json:"date"`
	EatenEnergyKcal int    `json:"eatenEnergyKcal"`
	ScanCount       int    `json:"scanCount"`
}

type weeklyEnergyTrendResponse struct {
	WeekStart string                          `json:"weekStart"`
	WeekEnd   string                          `json:"weekEnd"`
	Days      []dailyEnergyTrendPointResponse `json:"days"`
}

func buildWeeklyEnergyTrend(weekStart time.Time, dbSummaries []DailyTrendSummary) weeklyEnergyTrendResponse {
	daysMap := make(map[string]DailyTrendSummary)
	for _, s := range dbSummaries {
		dateStr := s.Date.Format("2006-01-02")
		daysMap[dateStr] = s
	}

	var days []dailyEnergyTrendPointResponse
	for i := 0; i < 7; i++ {
		currentDay := weekStart.AddDate(0, 0, i)
		dateStr := currentDay.Format("2006-01-02")
		summary := daysMap[dateStr]
		
		days = append(days, dailyEnergyTrendPointResponse{
			Date:            dateStr,
			EatenEnergyKcal: summary.EatenEnergyKcal,
			ScanCount:       summary.ScanCount,
		})
	}

	weekEnd := weekStart.AddDate(0, 0, 6).Format("2006-01-02")

	return weeklyEnergyTrendResponse{
		WeekStart: weekStart.Format("2006-01-02"),
		WeekEnd:   weekEnd,
		Days:      days,
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
