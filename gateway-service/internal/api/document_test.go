package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
)

type stubDocService struct {
	view *service.DocumentView
	err  error
}

func (s *stubDocService) GetDocument(ctx context.Context, id string) (*service.DocumentView, error) {
	return s.view, s.err
}

func TestDocumentHandler(t *testing.T) {
	tests := []struct {
		name             string
		stub             stubDocService
		wantCode         int
		expectPrediction bool
	}{
		{"view with document and prediction", stubDocService{
			view: &service.DocumentView{
				Document: &model.Document{
					ID:               "1",
					OriginalFilename: "invoice.pdf",
					StoredFilename:   "123-123.pdf",
					Status:           model.StatusDone,
					CreatedAt:        time.Now().UTC(),
				},
				Prediction: &model.Prediction{
					ID:         "1",
					DocumentID: "1",
					Label:      "finance",
					Confidence: 0.8,
					CreatedAt:  time.Now().UTC(),
				},
			}}, 200, true},
		{"view with document and no prediction", stubDocService{
			view: &service.DocumentView{
				Document: &model.Document{
					ID:               "1",
					OriginalFilename: "invoice.pdf",
					StoredFilename:   "123-123.pdf",
					Status:           model.StatusDone,
					CreatedAt:        time.Now().UTC(),
				},
			}}, 200, false},
		{"no documents by id", stubDocService{
			err: pgx.ErrNoRows,
		}, 404, false},
		{"unexpected error", stubDocService{
			err: fmt.Errorf("some error"),
		}, 500, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get("/documents/{id}", DocumentHandler(&test.stub))
			req := httptest.NewRequest(http.MethodGet, "/documents/1", nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if test.wantCode != rec.Code {
				t.Errorf("got status code %d, want %d", rec.Code, test.wantCode)
			}

			body := rec.Body.String()

			if strings.Contains(body, "stored_filename") || strings.Contains(body, "123-123.pdf") {
				t.Fatalf("internal storage field leaked: %s", body)
			}

			if test.expectPrediction {
				if !strings.Contains(body, "prediction") {
					t.Fatalf("expected prediction key, got %s", body)
				}

				var docResponse documentResponse
				err := json.NewDecoder(rec.Body).Decode(&docResponse)
				if err != nil {
					t.Fatalf("expected decoding of documentResponse, got error %s", err)
				}

				if docResponse.Prediction == nil {
					t.Fatal("expected prediction")
				}

				if docResponse.Prediction.Label != test.stub.view.Prediction.Label {
					t.Fatalf("expected prediction label %s, got %s", test.stub.view.Prediction.Label, docResponse.Prediction.Label)
				}

				if docResponse.Prediction.Confidence != test.stub.view.Prediction.Confidence {
					t.Fatalf("expected prediction confidence %f, got %f", test.stub.view.Prediction.Confidence, docResponse.Prediction.Confidence)
				}

				return
			}
			if strings.Contains(body, "prediction") {
				t.Fatalf("expected no prediction key, got %s", body)
			}
		})
	}

}
