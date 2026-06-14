package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type PostgresDocumentRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDocumentRepository(db *pgxpool.Pool) *PostgresDocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

func (r *PostgresDocumentRepository) Save(ctx context.Context, doc *model.Document) error {
	_, err := r.db.Exec(ctx, "INSERT INTO documents (id, original_filename, stored_filename, status, created_at) VALUES ($1, $2, $3, $4, $5)",
		doc.ID, doc.OriginalFilename, doc.StoredFilename, doc.Status, doc.CreatedAt)
	return err
}
