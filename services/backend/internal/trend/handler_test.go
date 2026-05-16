package trend

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
)

type mockStore struct {
	summaries []DailyTrendSummary
	err       error
}

func (m *mockStore) GetWeeklyTrendSummaries(ctx context.Context, anonymousUserID string, weekStart, weekEnd time.Time) ([]DailyTrendSummary, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.summaries, nil
}

func TestGetWeeklyTrend(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	
	// mock current time to be Wednesday, May 13 2026
	clockTime := time.Date(2026, 5, 13, 12, 0, 0, 0, time.Local)
	
	// Monday of that week is May 11
	mondayDate := time.Date(2026, 5, 11, 0, 0, 0, 0, time.Local)
	tuesdayDate := time.Date(2026, 5, 12, 0, 0, 0, 0, time.Local)

	store := &mockStore{
		summaries: []DailyTrendSummary{
			{Date: mondayDate, EatenEnergyKcal: 2100, ScanCount: 3},
			{Date: tuesdayDate, EatenEnergyKcal: 1900, ScanCount: 4},
		},
	}
	handler := NewHandler(store, nil)
	handler.clock = func() time.Time {
		return clockTime
	}

	router := chi.NewRouter()
	handler.RegisterRoutes(router, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.WithTestAnonymousUser(r.Context(), anonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/trends/weekly", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response weeklyEnergyTrendResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.WeekStart != "2026-05-11" {
		t.Fatalf("expected weekStart %q, got %q", "2026-05-11", response.WeekStart)
	}
	if response.WeekEnd != "2026-05-17" { // 6 days after Monday
		t.Fatalf("expected weekEnd %q, got %q", "2026-05-17", response.WeekEnd)
	}
	if len(response.Days) != 7 {
		t.Fatalf("expected 7 days, got %d", len(response.Days))
	}
	
	monday := response.Days[0]
	if monday.Date != "2026-05-11" || monday.EatenEnergyKcal != 2100 || monday.ScanCount != 3 {
		t.Fatalf("unexpected monday summary: %+v", monday)
	}

	tuesday := response.Days[1]
	if tuesday.Date != "2026-05-12" || tuesday.EatenEnergyKcal != 1900 || tuesday.ScanCount != 4 {
		t.Fatalf("unexpected tuesday summary: %+v", tuesday)
	}

	wednesday := response.Days[2]
	if wednesday.Date != "2026-05-13" || wednesday.EatenEnergyKcal != 0 || wednesday.ScanCount != 0 {
		t.Fatalf("unexpected wednesday summary: %+v", wednesday)
	}
}
