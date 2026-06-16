package main

import (
	"context"
	"errors"
	"log/slog"
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
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	r := chi.NewRouter()
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		slog.Error("Database URL not set")
		os.Exit(1)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		slog.Error("Could not create upload dir", "error", err)
		os.Exit(1)
	}

	runMigrations(cfg.DatabaseURL)

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("Could not initialize connection pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	fileStorage := storage.NewLocalFileStorage(cfg.UploadDir)
	documentRepository := repository.NewPostgresDocumentRepository(pool)
	docSvc := service.NewDocumentService(fileStorage, documentRepository)

	r.Get("/health", api.HealthHandler)
	r.Post("/upload", api.UploadHandler(docSvc))
	r.Get("/documents/{id}", api.DocumentHandler(docSvc))

	slog.Info("Server started on port", "port", cfg.Port)
	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		slog.Error("Could not run the application", "error", err)
	}
}

func runMigrations(databaseURL string) {
	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	m, err := migrate.New("file://migrations", migrateURL)
	if err != nil {
		slog.Error("failed to create migrator", "error", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			slog.Error("failed to run migrations", "error", err)
			os.Exit(1)
		}
		slog.Info("Nothing to migrate.")
	} else {
		slog.Info("Migrations applied.")
	}
}
