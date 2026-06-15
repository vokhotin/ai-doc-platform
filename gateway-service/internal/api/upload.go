package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
)

type documentService interface {
	Upload(ctx context.Context, file multipart.File, filename string) (*service.UploadResult, error)
}

func UploadHandler(svc documentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "file is required", http.StatusBadRequest)
			return
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				slog.Warn("failed to close upload file", "error", err)
			}
		}(file)

		safeName := filepath.Base(header.Filename)
		if safeName == "." || safeName == "/" {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		result, err := svc.Upload(
			r.Context(),
			file,
			safeName,
		)
		if err != nil {
			slog.Error("failed to upload file", "error", err)
			http.Error(w, "failed to upload file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(map[string]string{
			"id":       result.ID,
			"filename": result.Filename,
		})
		if err != nil {
			slog.Error("failed to encode response", "error", err)
			return
		}
	}
}
