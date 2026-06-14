package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/api"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/config"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/repository"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/service"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/storage"
)

func main() {
	r := chi.NewRouter()
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("Database URL not set")
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal(err)
	}

	runMigrations(cfg.DatabaseURL)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	fileStorage := storage.NewLocalFileStorage(cfg.UploadDir)
	documentRepository := repository.NewPostgresDocumentRepository(pool)
	docSvc := service.NewDocumentService(fileStorage, documentRepository)

	r.Get("/health", api.HealthHandler)
	r.Post("/upload", api.UploadHandler(docSvc))
	r.Get("/documents/{id}", api.DocumentHandler(docSvc))

	log.Printf("Server started on port %v", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal(err)
	}
}

func runMigrations(databaseURL string) {
	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	m, err := migrate.New("file://migrations", migrateURL)
	if err != nil {
		log.Fatal("failed to create migrator: ", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal("failed to run migrations: ", err)
		}
		log.Printf("Nothing to migrate.")
	} else {
		log.Printf("Migrations applied.")
	}
}
