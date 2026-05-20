package nudge

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func TestRecordNudgeResponseSucceeds(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: uuid.NewString()}
	nudgeID := uuid.NewString()
	scanID := uuid.NewString()

	store := &fakeNudgeStore{}
	router := newNudgeRouter(t, store, anonymousUser)

	body, _ := json.Marshal(map[string]string{
		"scanId":   scanID,
		"response": "followed",
	})
	req := httptest.NewRequest(http.MethodPost, "/nudges/"+nudgeID+"/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response["nudgeId"] != nudgeID {
		t.Fatalf("expected nudgeId %q, got %q", nudgeID, response["nudgeId"])
	}
	if response["scanId"] != scanID {
		t.Fatalf("expected scanId %q, got %q", scanID, response["scanId"])
	}
	if response["response"] != "followed" {
		t.Fatalf("expected response %q, got %q", "followed", response["response"])
	}
	if _, ok := response["nudgeResponseId"]; !ok {
		t.Fatal("expected nudgeResponseId in response")
	}
}

func TestRecordNudgeResponseRejectsMalformedJSON(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: uuid.NewString()}
	nudgeID := uuid.NewString()

	store := &fakeNudgeStore{}
	router := newNudgeRouter(t, store, anonymousUser)

	req := httptest.NewRequest(http.MethodPost, "/nudges/"+nudgeID+"/responses", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)
	if response["error"] != "invalid_nudge_response_request" {
		t.Fatalf("expected invalid_nudge_response_request error, got %q", response["error"])
	}
}

func TestRecordNudgeResponseRejectsInvalidResponseValue(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: uuid.NewString()}
	nudgeID := uuid.NewString()
	scanID := uuid.NewString()

	store := &fakeNudgeStore{}
	router := newNudgeRouter(t, store, anonymousUser)

	body, _ := json.Marshal(map[string]string{
		"scanId":   scanID,
		"response": "invalid_value",
	})
	req := httptest.NewRequest(http.MethodPost, "/nudges/"+nudgeID+"/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusBadRequest, rec.Code, rec.Body.String())
	}

	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)
	if response["error"] != "invalid_nudge_response_request" {
		t.Fatalf("expected invalid_nudge_response_request error, got %q", response["error"])
	}
}

func TestRecordNudgeResponseReturnsNotFoundForUnownedNudge(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: uuid.NewString()}
	nudgeID := uuid.NewString()
	scanID := uuid.NewString()

	store := &fakeNudgeStore{verifyErr: ErrNudgeNotFound}
	router := newNudgeRouter(t, store, anonymousUser)

	body, _ := json.Marshal(map[string]string{
		"scanId":   scanID,
		"response": "followed",
	})
	req := httptest.NewRequest(http.MethodPost, "/nudges/"+nudgeID+"/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}

	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)
	if response["error"] != "nudge_not_found" {
		t.Fatalf("expected nudge_not_found error, got %q", response["error"])
	}
}

func TestRecordNudgeResponseReturnsNotFoundForInvalidNudgeID(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: uuid.NewString()}
	scanID := uuid.NewString()

	store := &fakeNudgeStore{}
	router := newNudgeRouter(t, store, anonymousUser)

	body, _ := json.Marshal(map[string]string{
		"scanId":   scanID,
		"response": "followed",
	})
	req := httptest.NewRequest(http.MethodPost, "/nudges/not-a-uuid/responses", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}

	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)
	if response["error"] != "nudge_not_found" {
		t.Fatalf("expected nudge_not_found error, got %q", response["error"])
	}
}

// --- test helpers ---

func newNudgeRouter(t *testing.T, store Store, anonymousUser user.AnonymousUser) *chi.Mux {
	t.Helper()

	handler := NewHandler(store, slog.Default())
	router := chi.NewRouter()
	handler.RegisterRoutes(router, testRequireAnonymousUser(anonymousUser))
	return router
}

func testRequireAnonymousUser(anonymousUser user.AnonymousUser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.WithTestAnonymousUser(r.Context(), anonymousUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type fakeNudgeStore struct {
	verifyErr error
	recordErr error
}

func (s *fakeNudgeStore) VerifyNudgeOwnership(ctx context.Context, anonymousUserID string, scanID string, nudgeID string) error {
	return s.verifyErr
}

func (s *fakeNudgeStore) RecordResponse(ctx context.Context, record ResponseRecord) (ResponseRecord, error) {
	if s.recordErr != nil {
		return ResponseRecord{}, s.recordErr
	}
	record.CreatedAt = time.Now()
	return record, nil
}
