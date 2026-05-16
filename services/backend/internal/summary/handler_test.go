package summary

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
	summaries []MealSummary
	err       error
}

func (m *mockStore) GetDailyMealSummaries(ctx context.Context, anonymousUserID string, dateStart, dateEnd time.Time) ([]MealSummary, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.summaries, nil
}

func TestGetDailySummary(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := &mockStore{
		summaries: []MealSummary{
			{MealType: "lunch", EatenEnergyKcal: 500, ScanCount: 2},
			{MealType: "snack", EatenEnergyKcal: 200, ScanCount: 1},
		},
	}
	handler := NewHandler(store, nil)
	handler.clock = func() time.Time {
		return time.Date(2026, 5, 16, 12, 0, 0, 0, time.Local)
	}

	router := chi.NewRouter()
	handler.RegisterRoutes(router, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.WithTestAnonymousUser(r.Context(), anonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/summaries/daily", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response dailyEnergySummaryResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Date != "2026-05-16" {
		t.Fatalf("expected date %q, got %q", "2026-05-16", response.Date)
	}
	if response.EatenEnergyKcal != 700 {
		t.Fatalf("expected eaten energy 700, got %d", response.EatenEnergyKcal)
	}
	if response.DailyGoalEnergyKcal != 2000 {
		t.Fatalf("expected daily goal 2000, got %d", response.DailyGoalEnergyKcal)
	}
	if response.RemainingEnergyKcal != 1300 {
		t.Fatalf("expected remaining energy 1300, got %d", response.RemainingEnergyKcal)
	}
	if len(response.Meals) != 4 {
		t.Fatalf("expected 4 meal slots, got %d", len(response.Meals))
	}
	
	lunch := response.Meals[1]
	if lunch.MealType != "lunch" || lunch.EatenEnergyKcal != 500 || lunch.ScanCount != 2 {
		t.Fatalf("unexpected lunch summary: %+v", lunch)
	}
}

func TestGetDailySummaryWithExplicitDate(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := &mockStore{}
	handler := NewHandler(store, nil)

	router := chi.NewRouter()
	handler.RegisterRoutes(router, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.WithTestAnonymousUser(r.Context(), anonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/summaries/daily?date=2026-05-15", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response dailyEnergySummaryResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.Date != "2026-05-15" {
		t.Fatalf("expected date %q, got %q", "2026-05-15", response.Date)
	}
}

func TestGetDailySummaryInvalidDate(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := &mockStore{}
	handler := NewHandler(store, nil)

	router := chi.NewRouter()
	handler.RegisterRoutes(router, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.WithTestAnonymousUser(r.Context(), anonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/summaries/daily?date=invalid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}
