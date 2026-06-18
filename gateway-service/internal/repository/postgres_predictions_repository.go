package repository

import (
	"context"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

func (pr *PostgresRepository) SavePredictionAndMarkDocumentDone(ctx context.Context, prediction *model.Prediction) error {
	tx, err := pr.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO predictions (id, document_id, label, confidence, created_at) VALUES ($1, $2, $3, $4, $5)",
		prediction.ID, prediction.DocumentID, prediction.Label, prediction.Confidence, prediction.CreatedAt)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE documents SET status = $1 WHERE id = $2", model.StatusDone, prediction.DocumentID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return err
}

func (pr *PostgresRepository) GetLatestPredictionByDocumentId(ctx context.Context, documentId string) (*model.Prediction, error) {
	row := pr.db.QueryRow(ctx, "SELECT id, document_id, label, confidence, created_at FROM predictions WHERE document_id = $1 ORDER BY created_at DESC LIMIT 1", documentId)
	prediction := &model.Prediction{}
	err := row.Scan(&prediction.ID, &prediction.DocumentID, &prediction.Label, &prediction.Confidence, &prediction.CreatedAt)
	if err != nil {
		return nil, err
	}
	return prediction, nil
}
