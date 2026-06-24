//go:build integration

package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

var testPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16",
		tcpostgres.WithDatabase("gateway"),
		tcpostgres.WithUsername("gateway"),
		tcpostgres.WithPassword("gateway"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		fmt.Println("Start container:", err)
		os.Exit(1)
	}

	dbURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		fmt.Println("conn string:", err)
		os.Exit(1)
	}

	migrateURL := strings.Replace(dbURL, "postgres://", "pgx5://", 1)
	mig, err := migrate.New("file://../../migrations", migrateURL)
	if err != nil {
		fmt.Println("migrator:", err)
		os.Exit(1)
	}
	if err := mig.Up(); err != nil {
		fmt.Println("migrator up:", err)
		os.Exit(1)
	}

	testPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Println("pool:", err)
		os.Exit(1)
	}

	code := m.Run()

	testPool.Close()
	_ = container.Terminate(ctx)
	os.Exit(code)
}

func resetDB(t *testing.T) {
	t.Helper()
	_, err := testPool.Exec(context.Background(), "TRUNCATE predictions, documents CASCADE")
	if err != nil {
		t.Fatalf("reset db: %v", err)
	}
}

func newPendingDocument() *model.Document {
	id := uuid.New().String()
	return &model.Document{
		ID:               id,
		OriginalFilename: "report.pdf",
		StoredFilename:   id + ".pdf",
		Status:           model.StatusPending,
		CreatedAt:        time.Now().UTC(),
	}
}

func newPrediction(documentID string) *model.Prediction {
	return &model.Prediction{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		Label:      "finance",
		Confidence: 0.8,
		CreatedAt:  time.Now().UTC(),
	}
}
