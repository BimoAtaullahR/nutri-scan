package scan

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPInferenceClientDecodesRecognizerPayload(t *testing.T) {
	expected := InferenceResult{
		ModelVersion: "food-model-v0",
		FoodCategory: FoodCategoryConfidence{
			Slug:            "nasi_goreng",
			ConfidenceScore: 0.87,
		},
		Alternatives: []FoodCategoryConfidence{
			{Slug: "sate", ConfidenceScore: 0.42},
		},
		CoarsePortion: "medium",
		EstimatedEnergyRange: &EnergyRange{
			MinKcal: 400,
			MaxKcal: 500,
		},
		IsLowConfidence:     false,
		ConfidenceThreshold: 0.6,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if err := r.ParseMultipartForm(8 << 20); err != nil {
			t.Fatalf("parse multipart: %v", err)
		}
		file, header, err := r.FormFile("image")
		if err != nil {
			t.Fatalf("image field missing: %v", err)
		}
		defer file.Close()
		if header.Header.Get("Content-Type") != "image/png" {
			t.Fatalf("expected image/png, got %q", header.Header.Get("Content-Type"))
		}
		uploaded, err := io.ReadAll(file)
		if err != nil {
			t.Fatalf("read uploaded image: %v", err)
		}
		if string(uploaded) != "img-bytes" {
			t.Fatalf("expected image bytes to round-trip, got %q", string(uploaded))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expected); err != nil {
			t.Fatalf("encode response: %v", err)
		}
	}))
	defer server.Close()

	client := NewHTTPInferenceClient(server.URL)
	result, err := client.InferScan(context.Background(), ScanImage{
		ContentType: "image/png",
		Bytes:       []byte("img-bytes"),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.ModelVersion != expected.ModelVersion {
		t.Fatalf("expected model version %q, got %q", expected.ModelVersion, result.ModelVersion)
	}
	if result.FoodCategory.Slug != expected.FoodCategory.Slug {
		t.Fatalf("expected food category %q, got %q", expected.FoodCategory.Slug, result.FoodCategory.Slug)
	}
	if result.EstimatedEnergyRange == nil || result.EstimatedEnergyRange.MinKcal != 400 {
		t.Fatalf("expected estimated energy range to decode, got %#v", result.EstimatedEnergyRange)
	}
}

func TestHTTPInferenceClientReturnsErrorForNon2xx(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewHTTPInferenceClient(server.URL)
	_, err := client.InferScan(context.Background(), ScanImage{
		ContentType: "image/png",
		Bytes:       []byte("img-bytes"),
	})
	if err == nil {
		t.Fatal("expected error for non-2xx response")
	}
}

func TestHTTPInferenceClientReturnsErrorForMalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	client := NewHTTPInferenceClient(server.URL)
	_, err := client.InferScan(context.Background(), ScanImage{
		ContentType: "image/png",
		Bytes:       []byte("img-bytes"),
	})
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
}
