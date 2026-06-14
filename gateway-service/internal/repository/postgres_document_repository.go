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

func (r *PostgresDocumentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	row := r.db.QueryRow(ctx, "SELECT id, original_filename, stored_filename, status, created_at FROM documents WHERE id = $1", id)
	doc := &model.Document{}
	err := row.Scan(&doc.ID, &doc.OriginalFilename, &doc.StoredFilename, &doc.Status, &doc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
