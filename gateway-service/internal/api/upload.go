package api

import (
	"encoding/json"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
)

type documentService interface {
	Upload(file multipart.File, filename string) (*service.UploadResult, error)
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
				log.Printf("failed to close upload file: %v", err)
			}
		}(file)

		safeName := filepath.Base(header.Filename)
		if safeName == "." || safeName == "/" {
			http.Error(w, "invalid filename", http.StatusBadRequest)
			return
		}

		result, err := svc.Upload(
			file,
			safeName,
		)
		if err != nil {
			log.Printf("failed to upload file: %v", err)
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
			log.Printf("failed to encode response: %v", err)
			return
		}
	}
}
