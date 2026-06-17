package repository

import (
	"context"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

func (pr *PostgresRepository) SavePrediction(ctx context.Context, prediction *model.Prediction) error {
	_, err := pr.db.Exec(ctx, "INSERT INTO predictions (id, document_id, label, confidence, created_at) VALUES ($1, $2, $3, $4, $5)",
		prediction.Id, prediction.DocumentID, prediction.Label, prediction.Confidence, prediction.CreatedAt)
	return err
}

func (pr *PostgresRepository) GetPredictionsByDocumentId(ctx context.Context, documentId string) ([]*model.Prediction, error) {
	rows, err := pr.db.Query(ctx, "SELECT id, document_id, label, confidence, created_at FROM predictions WHERE document_id = $1", documentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var predictions []*model.Prediction
	for rows.Next() {
		prediction := &model.Prediction{}
		if err := rows.Scan(&prediction.Id, &prediction.DocumentID, &prediction.Label, &prediction.Confidence, &prediction.CreatedAt); err != nil {
			return predictions, err
		}
		predictions = append(predictions, prediction)
	}
	if err := rows.Err(); err != nil {
		return predictions, err
	}
	return predictions, nil
}
