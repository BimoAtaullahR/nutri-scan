package user

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestCreateAnonymousUserPersistsTokenHash(t *testing.T) {
	store := &fakeStore{}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/anonymous-users", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response struct {
		AnonymousUserID string `json:"anonymousUserId"`
		AccessToken     string `json:"accessToken"`
		TokenType       string `json:"tokenType"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if response.AnonymousUserID == "" {
		t.Fatal("expected anonymous user id")
	}
	if response.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if response.TokenType != "Bearer" {
		t.Fatalf("expected Bearer token type, got %q", response.TokenType)
	}
	if store.created.ID != response.AnonymousUserID {
		t.Fatalf("expected persisted id %q, got %q", response.AnonymousUserID, store.created.ID)
	}
	if store.created.TokenHash == "" {
		t.Fatal("expected persisted token hash")
	}
	if store.created.TokenHash == response.AccessToken {
		t.Fatal("expected persisted token hash, got raw token")
	}
	if store.created.TokenHash != hashBearerToken(response.AccessToken) {
		t.Fatal("expected persisted token hash to match returned bearer token")
	}
}

func TestRequireAnonymousUserRejectsMissingBearerToken(t *testing.T) {
	store := &fakeStore{}
	handler := NewHandler(store, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()

	handler.RequireAnonymousUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("protected handler should not run")
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestRequireAnonymousUserLoadsUserFromBearerToken(t *testing.T) {
	token := "nutriscan_test"
	expectedUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			expectedUser.TokenHash: expectedUser,
		},
	}
	handler := NewHandler(store, slog.Default())

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler.RequireAnonymousUser(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		anonymousUser, ok := AnonymousUserFromContext(r.Context())
		if !ok {
			t.Fatal("expected anonymous user in context")
		}
		if anonymousUser.ID != expectedUser.ID {
			t.Fatalf("expected anonymous user %q, got %q", expectedUser.ID, anonymousUser.ID)
		}
		w.WriteHeader(http.StatusNoContent)
	})).ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}
}

func TestUpdateProfilePersistsBMICategoryForAnonymousUser(t *testing.T) {
	token := "nutriscan_profile"
	anonymousUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			anonymousUser.TokenHash: anonymousUser,
		},
	}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPut, "/me/profile", strings.NewReader(`{
		"heightCm": 170,
		"weightKg": 72,
		"ageRange": "25_34"
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if store.upsertedProfile.AnonymousUserID != anonymousUser.ID {
		t.Fatalf("expected profile owner %q, got %q", anonymousUser.ID, store.upsertedProfile.AnonymousUserID)
	}
	if store.upsertedProfile.BMI != 24.91 {
		t.Fatalf("expected BMI 24.91, got %.2f", store.upsertedProfile.BMI)
	}
	if store.upsertedProfile.BMICategory != "normal" {
		t.Fatalf("expected normal BMI category, got %q", store.upsertedProfile.BMICategory)
	}

	var response userProfileResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.BMI != 24.91 {
		t.Fatalf("expected response BMI 24.91, got %.2f", response.BMI)
	}
	if response.BMICategory != "normal" {
		t.Fatalf("expected response BMI category normal, got %q", response.BMICategory)
	}
}

func TestGetProfileReturnsCurrentAnonymousUserProfile(t *testing.T) {
	token := "nutriscan_profile"
	anonymousUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			anonymousUser.TokenHash: anonymousUser,
		},
		profilesByUserID: map[string]UserProfile{
			anonymousUser.ID: {
				AnonymousUserID: anonymousUser.ID,
				HeightCm:        160,
				WeightKg:        80,
				BMI:             31.25,
				BMICategory:     "obese",
			},
		},
	}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d with body %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var response userProfileResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.BMICategory != "obese" {
		t.Fatalf("expected obese BMI category, got %q", response.BMICategory)
	}
}

func TestGetProfileReturnsNotFoundWhenProfileIsMissing(t *testing.T) {
	token := "nutriscan_profile"
	anonymousUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			anonymousUser.TokenHash: anonymousUser,
		},
	}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/me/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestUpdateProfileRejectsInvalidProfileRequest(t *testing.T) {
	token := "nutriscan_profile"
	anonymousUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			anonymousUser.TokenHash: anonymousUser,
		},
	}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPut, "/me/profile", strings.NewReader(`{
		"heightCm": 170,
		"weightKg": 72,
		"ageRange": "unknown"
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestUpdateProfileRejectsOutOfRangeMeasurements(t *testing.T) {
	token := "nutriscan_profile"
	anonymousUser := AnonymousUser{
		ID:        "9fd2e4f6-a1ab-432a-86bc-30a743a6f74b",
		TokenHash: hashBearerToken(token),
	}
	store := &fakeStore{
		usersByTokenHash: map[string]AnonymousUser{
			anonymousUser.TokenHash: anonymousUser,
		},
	}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPut, "/me/profile", strings.NewReader(`{
		"heightCm": 301,
		"weightKg": 72
	}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
	if store.upsertedProfile.AnonymousUserID != "" {
		t.Fatal("expected invalid User Profile request not to be persisted")
	}
}

func TestBMICategoryBoundaries(t *testing.T) {
	tests := map[float64]string{
		18.49: "underweight",
		18.5:  "normal",
		24.99: "normal",
		25:    "overweight",
		29.99: "overweight",
		30:    "obese",
	}

	for bmi, expected := range tests {
		if got := bmiCategory(bmi); got != expected {
			t.Fatalf("expected BMI %.2f to be %q, got %q", bmi, expected, got)
		}
	}
}

type fakeStore struct {
	created          AnonymousUser
	createErr        error
	usersByTokenHash map[string]AnonymousUser
	profilesByUserID map[string]UserProfile
	upsertedProfile  UserProfile
	profileErr       error
	upsertErr        error
}

func (s *fakeStore) CreateAnonymousUser(ctx context.Context, anonymousUser AnonymousUser) error {
	s.created = anonymousUser
	return s.createErr
}

func (s *fakeStore) GetAnonymousUserByTokenHash(ctx context.Context, tokenHash string) (AnonymousUser, error) {
	anonymousUser, ok := s.usersByTokenHash[tokenHash]
	if !ok {
		return AnonymousUser{}, ErrAnonymousUserNotFound
	}

	return anonymousUser, nil
}

func (s *fakeStore) GetUserProfile(ctx context.Context, anonymousUserID string) (UserProfile, error) {
	if s.profileErr != nil {
		return UserProfile{}, s.profileErr
	}

	profile, ok := s.profilesByUserID[anonymousUserID]
	if !ok {
		return UserProfile{}, ErrUserProfileNotFound
	}

	return profile, nil
}

func (s *fakeStore) UpsertUserProfile(ctx context.Context, profile UserProfile) (UserProfile, error) {
	if s.upsertErr != nil {
		return UserProfile{}, s.upsertErr
	}

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = profile.CreatedAt
	s.upsertedProfile = profile
	return profile, nil
}

func TestCreateAnonymousUserHandlesStoreFailure(t *testing.T) {
	store := &fakeStore{createErr: errors.New("store failed")}
	router := chi.NewRouter()
	NewHandler(store, slog.Default()).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/anonymous-users", strings.NewReader(""))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}
