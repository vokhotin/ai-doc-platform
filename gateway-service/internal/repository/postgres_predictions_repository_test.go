//go:build integration

package repository

import (
	"testing"

	"github.com/google/uuid"
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
