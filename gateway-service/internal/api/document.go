package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
)

type predictionResponse struct {
	Label      string    `json:"label"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}

type documentResponse struct {
	ID         string              `json:"id"`
	Filename   string              `json:"filename"`
	Status     string              `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	Prediction *predictionResponse `json:"prediction,omitempty"`
}

type documentGetService interface {
	GetDocument(ctx context.Context, id string) (*service.DocumentView, error)
}

func DocumentHandler(svc documentGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentID := chi.URLParam(r, "id")
		documentView, err := svc.GetDocument(r.Context(), documentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Info("Document does not exist", "id", documentID)
				http.Error(w, "no document found", http.StatusNotFound)
				return
			}

			http.Error(w, "could not fetch the doc", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := &documentResponse{
			ID:        documentView.Document.ID,
			Filename:  documentView.Document.OriginalFilename,
			Status:    string(documentView.Document.Status),
			CreatedAt: documentView.Document.CreatedAt,
		}
		if documentView.Prediction != nil {
			resp.Prediction = &predictionResponse{
				Label:      documentView.Prediction.Label,
				Confidence: documentView.Prediction.Confidence,
				CreatedAt:  documentView.Prediction.CreatedAt,
			}
		}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			slog.Error("failed to encode response", "error", err)
			return
		}
	}
}
