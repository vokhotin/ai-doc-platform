package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/api"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/config"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/storage"
)

func main() {
	r := chi.NewRouter()
	cfg := config.Load()

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	fileStorage := storage.NewLocalFileStorage(cfg.UploadDir)
	docSvc := service.NewDocumentService(fileStorage)

	r.Get("/health", api.HealthHandler)
	r.Post("/upload", api.UploadHandler(docSvc))

	log.Printf("Server started on port %v", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}
