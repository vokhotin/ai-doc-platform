//go:build integration

package repository

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

func TestSavePredictionAndMarkDocumentDone_Commits(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	doc := newPendingDocument()
	require.NoError(t, repo.SaveDocument(ctx, doc))

	pred := newPrediction(doc.ID)
	require.NoError(t, repo.SavePredictionAndMarkDocumentDone(ctx, pred))

	var labelGot string
	var confGot float64
	require.NoError(t, testPool.QueryRow(ctx,
		`SELECT label, confidence FROM predictions WHERE document_id = $1`,
		doc.ID,
	).Scan(&labelGot, &confGot))
	assert.Equal(t, "finance", labelGot)
	assert.Equal(t, 0.8, confGot)

	var docStatus model.DocumentStatus
	require.NoError(t, testPool.QueryRow(ctx,
		`SELECT status FROM documents WHERE id = $1`,
		doc.ID,
	).Scan(&docStatus))

	assert.Equal(t, model.StatusDone, docStatus)
}

func TestSavePredictionAndMarkDocumentDone_RollsBackOnFKViolation(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	pred := newPrediction(uuid.New().String())
	err := repo.SavePredictionAndMarkDocumentDone(ctx, pred)
	require.Error(t, err)

	var count int
	require.NoError(t, testPool.QueryRow(ctx, "SELECT COUNT(*) FROM predictions").Scan(&count))
	assert.Equal(t, 0, count)
}

func TestPostgresRepository_GetLatestPredictionByDocumentId(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	doc := newPendingDocument()
	require.NoError(t, repo.SaveDocument(ctx, doc))

	docID := doc.ID

	older := newPrediction(docID)
	older.CreatedAt = time.Date(
		2000, 11, 17, 20, 34, 58, 651387237, time.UTC)
	newer := newPrediction(docID)

	_, err := testPool.Exec(ctx, "INSERT INTO predictions (id, document_id, label, confidence, created_at) VALUES ($1, $2, $3, $4, $5)",
		older.ID, older.DocumentID, older.Label, older.Confidence, older.CreatedAt)
	require.NoError(t, err)
	_, err = testPool.Exec(ctx, "INSERT INTO predictions (id, document_id, label, confidence, created_at) VALUES ($1, $2, $3, $4, $5)",
		newer.ID, newer.DocumentID, newer.Label, newer.Confidence, newer.CreatedAt)
	require.NoError(t, err)

	prediction, err := repo.GetLatestPredictionByDocumentId(ctx, docID)
	require.NoError(t, err)

	assert.Equal(t, newer.ID, prediction.ID)
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	_, err := repo.GetByID(ctx, uuid.New().String())
	require.ErrorIs(t, err, pgx.ErrNoRows)
}

func TestPostgresRepository_GetLatestPredictionByDocumentId_NotFound(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	_, err := repo.GetLatestPredictionByDocumentId(ctx, uuid.New().String())
	require.ErrorIs(t, err, pgx.ErrNoRows)
}
