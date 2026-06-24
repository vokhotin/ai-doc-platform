//go:build integration

package repository

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

func TestPostgresRepository_SaveDocumentAndReadItBack(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)

	doc := newPendingDocument()
	require.NoError(t, repo.SaveDocument(ctx, doc))
	document, err := repo.GetByID(ctx, doc.ID)
	require.NoError(t, err)

	assert.Equal(t, doc.ID, document.ID)
	assert.Equal(t, doc.OriginalFilename, document.OriginalFilename)
	assert.Equal(t, doc.StoredFilename, document.StoredFilename)
	assert.Equal(t, doc.Status, document.Status)
	assert.WithinDuration(t, doc.CreatedAt, document.CreatedAt, time.Second)
}

func TestPostgresRepository_UpdateDocumentStatus(t *testing.T) {
	resetDB(t)
	ctx := t.Context()
	repo := NewPostgresRepository(testPool)
	doc := newPendingDocument()
	require.NoError(t, repo.SaveDocument(ctx, doc))

	require.NoError(t, repo.UpdateDocumentStatus(ctx, doc.ID, model.StatusDone))

	var docStatus model.DocumentStatus
	require.NoError(t, testPool.QueryRow(ctx,
		`SELECT status FROM documents WHERE id = $1`,
		doc.ID,
	).Scan(&docStatus))
	assert.Equal(t, model.StatusDone, docStatus)
}
