package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"time"
)

type ScanImage struct {
	ContentType string
	Bytes       []byte
}

type FoodCategoryConfidence struct {
	Slug            string  `json:"slug"`
	ConfidenceScore float64 `json:"confidenceScore"`
}

type InferenceResult struct {
	ModelVersion         string                   `json:"modelVersion"`
	FoodCategory         FoodCategoryConfidence   `json:"foodCategory"`
	Alternatives         []FoodCategoryConfidence `json:"alternatives"`
	CoarsePortion        string                   `json:"coarsePortion"`
	EstimatedEnergyRange *EnergyRange             `json:"estimatedEnergyRange"`
	IsLowConfidence      bool                     `json:"isLowConfidence"`
	ConfidenceThreshold  float64                  `json:"confidenceThreshold"`
}

type InferenceClient interface {
	InferScan(ctx context.Context, image ScanImage) (InferenceResult, error)
}

type HTTPInferenceClient struct {
	endpoint   string
	httpClient *http.Client
}

func NewHTTPInferenceClient(baseURL string) *HTTPInferenceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	endpoint := baseURL + "/infer"
	if parsed, err := url.Parse(baseURL); err == nil && strings.HasSuffix(parsed.Path, "/infer") {
		endpoint = baseURL
	}

	return &HTTPInferenceClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *HTTPInferenceClient) InferScan(ctx context.Context, image ScanImage) (InferenceResult, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreatePart(textproto.MIMEHeader{
		"Content-Disposition": {`form-data; name="image"; filename="scan-image"`},
		"Content-Type":        {image.ContentType},
	})
	if err != nil {
		return InferenceResult{}, err
	}
	if _, err := part.Write(image.Bytes); err != nil {
		return InferenceResult{}, err
	}
	if err := writer.Close(); err != nil {
		return InferenceResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, &body)
	if err != nil {
		return InferenceResult{}, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return InferenceResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		_, _ = io.Copy(io.Discard, resp.Body)
		return InferenceResult{}, fmt.Errorf("ai inference returned status %d", resp.StatusCode)
	}

	var result InferenceResult
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return InferenceResult{}, err
	}

	return result, nil
}
