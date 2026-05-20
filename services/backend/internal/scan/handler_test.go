package scan

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/nudge"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/user"
	"github.com/go-chi/chi/v5"
)

func TestCreateScanCompletesSyncFirstWithFakeAIClient(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{
		result: validInferenceResult(),
	}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "lunch", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if store.created.AnonymousUserID != anonymousUser.ID {
		t.Fatalf("expected scan owner %q, got %q", anonymousUser.ID, store.created.AnonymousUserID)
	}
	if store.created.MealType != MealTypeLunch {
		t.Fatalf("expected meal type %q, got %q", MealTypeLunch, store.created.MealType)
	}
	if len(aiClient.receivedImage.Bytes) != len("not-empty") {
		t.Fatal("expected Scan Image bytes to be forwarded to AI/ML Inference client")
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusCompleted {
		t.Fatalf("expected completed status, got %q", response.Status)
	}
	if response.EstimatedEnergyKcal == nil || *response.EstimatedEnergyKcal != 450 {
		t.Fatalf("expected estimated energy midpoint 450, got %#v", response.EstimatedEnergyKcal)
	}
	if response.Inference == nil || response.Inference.FoodCategory != "nasi_goreng" {
		t.Fatalf("expected inference summary, got %#v", response.Inference)
	}
	if response.NudgeDecision == nil {
		t.Fatal("expected backend-owned Nudge Decision")
	}
	if response.NudgeDecision.Action != nudge.ActionEatAsPlanned {
		t.Fatalf("expected eat-as-planned Nudge Action, got %q", response.NudgeDecision.Action)
	}
	if response.NudgeDecision.Basis != nudge.BasisGeneric {
		t.Fatalf("expected generic Nudge Decision basis, got %q", response.NudgeDecision.Basis)
	}
	if response.NudgeDecision.IsPersonalized {
		t.Fatal("expected Generic Nudge Decision when User Profile is missing")
	}
}

func TestCreateScanProducesReviewFoodNudgeForLowConfidenceInference(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	result := validInferenceResult()
	result.IsLowConfidence = true
	result.FoodCategory.Slug = "unknown_food"
	result.FoodCategory.ConfidenceScore = 0.42
	result.Alternatives = []FoodCategoryConfidence{
		{Slug: "sate", ConfidenceScore: 0.42},
		{Slug: "rendang", ConfidenceScore: 0.31},
	}
	result.EstimatedEnergyRange = nil
	aiClient := &fakeInferenceClient{result: result}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "lunch", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusCompleted {
		t.Fatalf("expected low-confidence Scan to complete, got %q", response.Status)
	}
	if response.FailureReason != nil {
		t.Fatalf("expected low-confidence Scan not to fail, got %q", *response.FailureReason)
	}
	if response.EstimatedEnergyKcal != nil {
		t.Fatalf("expected low-confidence Scan not to estimate energy, got %d", *response.EstimatedEnergyKcal)
	}
	if response.Inference == nil || response.Inference.FoodCategory != "unknown_food" {
		t.Fatalf("expected Unknown Food inference summary, got %#v", response.Inference)
	}
	if len(response.Inference.Alternatives) != 2 {
		t.Fatalf("expected review alternatives, got %#v", response.Inference.Alternatives)
	}
	if response.NudgeDecision == nil {
		t.Fatal("expected Review Food Nudge")
	}
	if response.NudgeDecision.Action != nudge.ActionReviewFood {
		t.Fatalf("expected review-food Nudge Action, got %q", response.NudgeDecision.Action)
	}
	if response.NudgeDecision.Basis != nudge.BasisReviewFood {
		t.Fatalf("expected review-food basis, got %q", response.NudgeDecision.Basis)
	}
}

func TestCreateScanProducesPersonalizedNudgeWhenProfileExists(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	result := validInferenceResult()
	result.EstimatedEnergyRange = &EnergyRange{MinKcal: 500, MaxKcal: 520}
	aiClient := &fakeInferenceClient{result: result}
	userStore := newFakeUserStore(anonymousUser)
	userStore.profilesByUserID = map[string]user.UserProfile{
		anonymousUser.ID: {
			AnonymousUserID: anonymousUser.ID,
			BMICategory:     "overweight",
		},
	}
	router := newAuthenticatedScanRouterWithUserStore(t, store, aiClient, userStore)

	req := newMultipartScanRequest(t, "/scans", "dinner", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.NudgeDecision == nil {
		t.Fatal("expected Personalized Nudge Decision")
	}
	if !response.NudgeDecision.IsPersonalized {
		t.Fatal("expected Nudge Decision to be personalized")
	}
	if response.NudgeDecision.Basis != nudge.BasisPersonalized {
		t.Fatalf("expected personalized basis, got %q", response.NudgeDecision.Basis)
	}
	if response.NudgeDecision.Action != nudge.ActionSetAsidePortion {
		t.Fatalf("expected set-aside-portion Nudge Action, got %q", response.NudgeDecision.Action)
	}
	if response.NudgeDecision.EstimatedPreventedEnergyKcal == nil {
		t.Fatal("expected estimated prevented energy for set-aside Nudge Action")
	}
}

func TestCreateScanMarksNudgeDecisionFailureFailed(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{result: validInferenceResult()}
	userStore := newFakeUserStore(anonymousUser)
	userStore.profileErr = errors.New("profile database unavailable")
	router := newAuthenticatedScanRouterWithUserStore(t, store, aiClient, userStore)

	req := newMultipartScanRequest(t, "/scans", "dinner", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusFailed {
		t.Fatalf("expected failed status, got %q", response.Status)
	}
	if response.FailureReason == nil || *response.FailureReason != "nudge_decision_failed" {
		t.Fatalf("expected nudge_decision_failed reason, got %#v", response.FailureReason)
	}
}

func TestCreateScanReturnsProcessingWhenInferenceExceedsSyncWindow(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{err: context.DeadlineExceeded}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "breakfast", []byte("not-empty"), "image/jpeg")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusAccepted, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusProcessing {
		t.Fatalf("expected processing status, got %q", response.Status)
	}
	if response.EstimatedEnergyKcal != nil || response.Inference != nil {
		t.Fatalf("expected no completed feedback for processing Scan, got %#v", response)
	}
}

func TestCreateScanMarksTechnicalInferenceFailureFailed(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{err: errors.New("ai unavailable")}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "snack", []byte("not-empty"), "image/webp")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusFailed {
		t.Fatalf("expected failed status, got %q", response.Status)
	}
	if response.FailureReason == nil || *response.FailureReason != "ai_inference_failed" {
		t.Fatalf("expected ai_inference_failed reason, got %#v", response.FailureReason)
	}
}

func TestCreateScanMarksInvalidInferencePayloadFailed(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{
		result: InferenceResult{
			ModelVersion: "food-model-v0",
			FoodCategory: FoodCategoryConfidence{
				Slug:            "nasi_goreng",
				ConfidenceScore: 1.2,
			},
			CoarsePortion:       "medium",
			ConfidenceThreshold: 0.6,
		},
	}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "snack", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var response scanResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.Status != ScanStatusFailed {
		t.Fatalf("expected failed status, got %q", response.Status)
	}
	if response.FailureReason == nil || *response.FailureReason != "invalid_inference_payload" {
		t.Fatalf("expected invalid_inference_payload reason, got %#v", response.FailureReason)
	}
}

func TestCreateScanAssignsFallbackMealType(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{result: validInferenceResult()}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)
	handler := router.scanHandler
	handler.clock = func() time.Time {
		return time.Date(2026, 5, 16, 18, 0, 0, 0, time.Local)
	}

	req := newMultipartScanRequest(t, "/scans", "", []byte("not-empty"), "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	if store.created.MealType != MealTypeDinner {
		t.Fatalf("expected fallback meal type %q, got %q", MealTypeDinner, store.created.MealType)
	}
}

func TestCreateScanRejectsInvalidUploadBeforeInference(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	store := newFakeStore()
	aiClient := &fakeInferenceClient{result: validInferenceResult()}
	router := newAuthenticatedScanRouter(t, store, aiClient, anonymousUser)

	req := newMultipartScanRequest(t, "/scans", "snack", nil, "image/png")
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if aiClient.calls != 0 {
		t.Fatal("expected invalid Scan Image not to be sent to AI/ML Inference")
	}
	if store.created.ID != "" {
		t.Fatal("expected invalid Scan Image not to create a Scan")
	}
}

func TestGetScanRequiresCurrentAnonymousUserOwnership(t *testing.T) {
	anonymousUser := user.AnonymousUser{ID: "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b"}
	otherUserID := "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	store := newFakeStore()
	scan := Scan{
		ID:              "33c31b25-0eba-48e5-b2ce-52eb49e84cc7",
		AnonymousUserID: otherUserID,
		Status:          ScanStatusCompleted,
		MealType:        MealTypeSnack,
	}
	store.scansByID[scan.ID] = scan
	router := newAuthenticatedScanRouter(t, store, &fakeInferenceClient{}, anonymousUser)

	req := httptest.NewRequest(http.MethodGet, "/scans/"+scan.ID, nil)
	req.Header.Set("Authorization", "Bearer token")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

type scanRouter struct {
	router      *chi.Mux
	scanHandler *Handler
}

func (r scanRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func newAuthenticatedScanRouter(t *testing.T, store *fakeScanStore, aiClient *fakeInferenceClient, anonymousUser user.AnonymousUser) scanRouter {
	t.Helper()

	return newAuthenticatedScanRouterWithUserStore(t, store, aiClient, newFakeUserStore(anonymousUser))
}

func newAuthenticatedScanRouterWithUserStore(t *testing.T, store *fakeScanStore, aiClient *fakeInferenceClient, userStore *fakeUserStore) scanRouter {
	t.Helper()

	userHandler := user.NewHandler(userStore, slog.Default())

	router := chi.NewRouter()
	scanHandler := NewHandler(store, userStore, aiClient, slog.Default())
	scanHandler.RegisterRoutes(router, userHandler.RequireAnonymousUser)

	return scanRouter{router: router, scanHandler: scanHandler}
}

func newMultipartScanRequest(t *testing.T, target string, mealType string, imageBytes []byte, contentType string) *http.Request {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	if mealType != "" {
		if err := writer.WriteField("mealType", mealType); err != nil {
			t.Fatalf("write meal type field: %v", err)
		}
	}

	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="image"; filename="scan-image"`},
		"Content-Type":        {contentType},
	})
	if err != nil {
		t.Fatalf("create image field: %v", err)
	}
	if _, err := part.Write(imageBytes); err != nil {
		t.Fatalf("write image field: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, target, &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func validInferenceResult() InferenceResult {
	return InferenceResult{
		ModelVersion: "food-model-v0",
		FoodCategory: FoodCategoryConfidence{
			Slug:            "nasi_goreng",
			ConfidenceScore: 0.87,
		},
		Alternatives: []FoodCategoryConfidence{
			{Slug: "fried_rice", ConfidenceScore: 0.61},
		},
		CoarsePortion: "medium",
		EstimatedEnergyRange: &EnergyRange{
			MinKcal: 400,
			MaxKcal: 500,
		},
		IsLowConfidence:     false,
		ConfidenceThreshold: 0.6,
	}
}

type fakeInferenceClient struct {
	result        InferenceResult
	err           error
	calls         int
	receivedImage ScanImage
}

func (c *fakeInferenceClient) InferScan(ctx context.Context, image ScanImage) (InferenceResult, error) {
	c.calls++
	c.receivedImage = image
	if c.err != nil {
		return InferenceResult{}, c.err
	}

	return c.result, nil
}

type fakeScanStore struct {
	created   Scan
	scansByID map[string]Scan
}

func newFakeStore() *fakeScanStore {
	return &fakeScanStore{
		scansByID: map[string]Scan{},
	}
}

func (s *fakeScanStore) CreateProcessingScan(ctx context.Context, scan Scan) error {
	now := time.Date(2026, 5, 16, 10, 0, 0, 0, time.UTC)
	scan.Status = ScanStatusProcessing
	scan.CreatedAt = now
	scan.UpdatedAt = now
	s.created = scan
	s.scansByID[scan.ID] = scan
	return nil
}

func (s *fakeScanStore) CompleteScan(ctx context.Context, anonymousUserID string, scanID string, inference InferenceResult, nudgeDecision nudge.Decision) (Scan, error) {
	scan, err := s.ownedScan(anonymousUserID, scanID)
	if err != nil {
		return Scan{}, err
	}
	scan.Status = ScanStatusCompleted
	scan.Inference = &inference
	scan.EstimatedEnergyRange = inference.EstimatedEnergyRange
	scan.NudgeDecision = &nudgeDecision
	scan.UpdatedAt = scan.UpdatedAt.Add(time.Second)
	s.scansByID[scanID] = scan
	return scan, nil
}

func (s *fakeScanStore) FailScan(ctx context.Context, anonymousUserID string, scanID string, reason string) (Scan, error) {
	scan, err := s.ownedScan(anonymousUserID, scanID)
	if err != nil {
		return Scan{}, err
	}
	scan.Status = ScanStatusFailed
	scan.FailureReason = reason
	scan.UpdatedAt = scan.UpdatedAt.Add(time.Second)
	s.scansByID[scanID] = scan
	return scan, nil
}

func (s *fakeScanStore) GetScan(ctx context.Context, anonymousUserID string, scanID string) (Scan, error) {
	return s.ownedScan(anonymousUserID, scanID)
}

func (s *fakeScanStore) ownedScan(anonymousUserID string, scanID string) (Scan, error) {
	scan, ok := s.scansByID[scanID]
	if !ok || scan.AnonymousUserID != anonymousUserID {
		return Scan{}, ErrScanNotFound
	}

	return scan, nil
}

type fakeUserStore struct {
	usersByTokenHash map[string]user.AnonymousUser
	profilesByUserID map[string]user.UserProfile
	profileErr       error
}

func newFakeUserStore(anonymousUser user.AnonymousUser) *fakeUserStore {
	tokenHash := hashTestToken("token")
	anonymousUser.TokenHash = tokenHash
	return &fakeUserStore{
		usersByTokenHash: map[string]user.AnonymousUser{
			tokenHash: anonymousUser,
		},
		profilesByUserID: map[string]user.UserProfile{},
	}
}

func (s *fakeUserStore) CreateAnonymousUser(ctx context.Context, anonymousUser user.AnonymousUser) error {
	return nil
}

func (s *fakeUserStore) GetAnonymousUserByTokenHash(ctx context.Context, tokenHash string) (user.AnonymousUser, error) {
	anonymousUser, ok := s.usersByTokenHash[tokenHash]
	if !ok {
		return user.AnonymousUser{}, user.ErrAnonymousUserNotFound
	}

	return anonymousUser, nil
}

func (s *fakeUserStore) GetUserProfile(ctx context.Context, anonymousUserID string) (user.UserProfile, error) {
	if s.profileErr != nil {
		return user.UserProfile{}, s.profileErr
	}

	profile, ok := s.profilesByUserID[anonymousUserID]
	if !ok {
		return user.UserProfile{}, user.ErrUserProfileNotFound
	}

	return profile, nil
}

func (s *fakeUserStore) UpsertUserProfile(ctx context.Context, profile user.UserProfile) (user.UserProfile, error) {
	return profile, nil
}

func hashTestToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
