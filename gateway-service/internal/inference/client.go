package inference

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type HTTPInferenceClient struct {
	InferenceURL string
	httpClient   *http.Client
}

type predictResponse struct {
	DocumentID string  `json:"document_id"`
	Label      string  `json:"label"`
	Confidence float64 `json:"confidence"`
}

type predictRequest struct {
	DocumentID string `json:"document_id"`
	Text       string `json:"text"`
}

func NewHTTPInferenceClient(httpInferenceURL string) *HTTPInferenceClient {
	return &HTTPInferenceClient{
		InferenceURL: httpInferenceURL,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *HTTPInferenceClient) Predict(ctx context.Context, documentID string, text string) (*model.InferenceResult, error) {
	body := predictRequest{
		DocumentID: documentID,
		Text:       text,
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.InferenceURL+"/predict", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("Could not call Inference service", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("Expected 200 OK response", "status", resp.Status, "request", body, "response", resp.Request.Body)
		return nil, fmt.Errorf("inference service responded with %s", resp.Status)
	}

	var result predictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("Could not decode response body", "error", err)
		return nil, err
	}

	return &model.InferenceResult{
		DocumentID: result.DocumentID,
		Label:      result.Label,
		Confidence: result.Confidence,
	}, nil
}
