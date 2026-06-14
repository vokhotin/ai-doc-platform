package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type documentGetService interface {
	GetDocument(ctx context.Context, id string) (*model.Document, error)
}

func DocumentHandler(svc documentGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentID := chi.URLParam(r, "id")
		document, err := svc.GetDocument(r.Context(), documentID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Printf("Document %v does not exist", documentID)
				http.Error(w, "no document found", http.StatusNotFound)
				return
			}

			http.Error(w, "could not fetch the doc", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(map[string]string{
			"id":              document.ID,
			"filename":        document.OriginalFilename,
			"stored_filename": document.StoredFilename,
			"status":          string(document.Status),
			"created_at":      document.CreatedAt.Format(time.RFC3339),
		})
		if err != nil {
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}
