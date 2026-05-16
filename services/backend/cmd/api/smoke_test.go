package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"
	"time"

	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/platform/config"
	"github.com/BimoAtaullahR/nutri-scan/services/backend/internal/platform/database"
	httpapi "github.com/BimoAtaullahR/nutri-scan/services/backend/internal/platform/http"
)

func TestCoreScanLoopSmoke(t *testing.T) {
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping end-to-end smoke test")
	}

	ctx := context.Background()
	db, err := database.Open(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	// Setup fake AI Inference Server
	aiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"modelVersion": "food-model-smoke",
			"foodCategory": {"slug": "nasi_goreng", "confidenceScore": 0.85},
			"alternatives": [],
			"coarsePortion": "medium",
			"estimatedEnergyRange": {"minKcal": 400, "maxKcal": 500},
			"isLowConfidence": false,
			"confidenceThreshold": 0.6
		}`))
	}))
	defer aiServer.Close()

	cfg := config.Config{
		HTTPAddr:       ":8080",
		DatabaseURL:    dbURL,
		AIInferenceURL: aiServer.URL,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	router := httpapi.NewRouter(cfg, db, logger)

	server := httptest.NewServer(router)
	defer server.Close()
	client := server.Client()

	// 1. Create Anonymous User
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/anonymous-users", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create user request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create user status: %d", resp.StatusCode)
	}

	var userResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&userResp)
	resp.Body.Close()

	token, ok := userResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("expected token, got: %v", userResp["token"])
	}

	// 2. Update User Profile
	profileBody := []byte(`{"heightCm": 170, "weightKg": 75, "ageRange": "20-29"}`)
	req, _ = http.NewRequest(http.MethodPut, server.URL+"/me/profile", bytes.NewReader(profileBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("update profile request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update profile status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// 3. Upload Scan
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("mealType", "lunch")
	part, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="image"; filename="scan.jpg"`},
		"Content-Type":        {"image/jpeg"},
	})
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ = http.NewRequest(http.MethodPost, server.URL+"/scans", &body)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("upload scan request failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		t.Fatalf("upload scan status: %d", resp.StatusCode)
	}

	var scanResp map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&scanResp)
	resp.Body.Close()

	scanID, _ := scanResp["scanId"].(string)

	// Poll if it is still processing
	for scanResp["status"] == "processing" {
		time.Sleep(100 * time.Millisecond)
		req, _ = http.NewRequest(http.MethodGet, server.URL+"/scans/"+scanID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, _ = client.Do(req)
		json.NewDecoder(resp.Body).Decode(&scanResp)
		resp.Body.Close()
	}

	if scanResp["status"] != "completed" {
		t.Fatalf("expected completed scan, got: %s", scanResp["status"])
	}

	// 4. Record Nudge Response
	nudgeDecision, ok := scanResp["nudgeDecision"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected nudge decision in scan response")
	}
	nudgeID, _ := nudgeDecision["nudgeId"].(string)

	nudgeRespBody := []byte(`{"scanId": "` + scanID + `", "response": "followed"}`)
	req, _ = http.NewRequest(http.MethodPost, server.URL+"/nudges/"+nudgeID+"/responses", bytes.NewReader(nudgeRespBody))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("record nudge response failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("record nudge status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// 5. Get Daily Summary
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/summaries/daily", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get daily summary failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get daily summary status: %d", resp.StatusCode)
	}
	resp.Body.Close()

	// 6. Get Weekly Trend
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/trends/weekly", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get weekly trend failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get weekly trend status: %d", resp.StatusCode)
	}
	resp.Body.Close()
}
