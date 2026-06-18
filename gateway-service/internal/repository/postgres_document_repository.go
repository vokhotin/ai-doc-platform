package repository

import (
	"context"
	"fmt"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

func (pr *PostgresRepository) SaveDocument(ctx context.Context, doc *model.Document) error {
	_, err := pr.db.Exec(ctx, "INSERT INTO documents (id, original_filename, stored_filename, status, created_at) VALUES ($1, $2, $3, $4, $5)",
		doc.ID, doc.OriginalFilename, doc.StoredFilename, doc.Status, doc.CreatedAt)
	return err
}

func (pr *PostgresRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	row := pr.db.QueryRow(ctx, "SELECT id, original_filename, stored_filename, status, created_at FROM documents WHERE id = $1", id)
	doc := &model.Document{}
	err := row.Scan(&doc.ID, &doc.OriginalFilename, &doc.StoredFilename, &doc.Status, &doc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (pr *PostgresRepository) UpdateDocumentStatus(ctx context.Context, documentID string, status model.DocumentStatus) error {
	exec, err := pr.db.Exec(ctx, "UPDATE documents SET status = $1 WHERE id = $2", status, documentID)
	if err != nil {
		return err
	}
	if exec.RowsAffected() == 0 {
		return fmt.Errorf("document %s not found", documentID)
	}
	return nil
}
