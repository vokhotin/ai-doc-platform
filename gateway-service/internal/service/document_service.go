package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type DocumentView struct {
	Document   *model.Document
	Prediction *model.Prediction
}

type FileStorage interface {
	Save(filename string, src io.Reader) error
}

type Repository interface {
	SaveDocument(ctx context.Context, doc *model.Document) error
	GetByID(ctx context.Context, id string) (*model.Document, error)
	UpdateDocumentStatus(ctx context.Context, documentID string, status model.DocumentStatus) error
	SavePredictionAndMarkDocumentDone(ctx context.Context, prediction *model.Prediction) error
	GetLatestPredictionByDocumentId(ctx context.Context, documentId string) (*model.Prediction, error)
}

type InferenceClient interface {
	Predict(ctx context.Context, documentID string, text string) (*model.InferenceResult, error)
}

type UploadResult struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type DocumentService struct {
	fs   FileStorage
	repo Repository
	ic   InferenceClient
}

func NewDocumentService(fs FileStorage, repo Repository, ic InferenceClient) *DocumentService {
	return &DocumentService{
		fs:   fs,
		repo: repo,
		ic:   ic,
	}
}

func (s *DocumentService) Upload(
	ctx context.Context,
	file multipart.File,
	filename string,
) (*UploadResult, error) {
	documentID := uuid.New().String()
	extension := filepath.Ext(filename)
	storedFilename := documentID + extension

	err := s.fs.Save(storedFilename, file)
	if err != nil {
		return nil, err
	}

	doc := &model.Document{
		ID:               documentID,
		OriginalFilename: filename,
		StoredFilename:   storedFilename,
		Status:           model.StatusPending,
		CreatedAt:        time.Now().UTC(),
	}
	err = s.repo.SaveDocument(ctx, doc)
	if err != nil {
		return nil, err
	}

	slog.Info("saved document", "id", doc.ID)

	inferenceResult, err := s.ic.Predict(ctx, doc.ID, doc.OriginalFilename)
	if err != nil {
		slog.Error("failed to predict type of document", "id", doc.ID, "error", err)
		errUpdate := s.updateDocumentStatus(ctx, documentID, model.StatusFailed)
		if errUpdate != nil {
			return nil, fmt.Errorf("failed to predict and failed to update document status. prediction error: \n%w,"+
				" \nUpdate document error\n%s", err, errUpdate)
		}
		return nil, err
	}

	slog.Info("predict result", "id", doc.ID, "label", inferenceResult.Label, "confidence", inferenceResult.Confidence)

	prediction := &model.Prediction{
		ID:         uuid.New().String(),
		DocumentID: inferenceResult.DocumentID,
		Label:      inferenceResult.Label,
		Confidence: inferenceResult.Confidence,
		CreatedAt:  time.Now().UTC(),
	}
	err = s.repo.SavePredictionAndMarkDocumentDone(ctx, prediction)
	if err != nil {
		return nil, err
	}

	slog.Info("save prediction and updated document", "id", doc.ID, "predictionID", prediction.ID)

	return &UploadResult{
		ID:       documentID,
		Filename: filename,
	}, nil
}

func (s *DocumentService) updateDocumentStatus(ctx context.Context, documentID string, status model.DocumentStatus) error {
	err := s.repo.UpdateDocumentStatus(ctx, documentID, status)
	if err != nil {
		slog.Error("failed to update document status", "id", documentID, "status", status, "error", err)
		return err
	}
	return nil
}

func (s *DocumentService) GetDocument(ctx context.Context, id string) (*DocumentView, error) {
	document, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	prediction, err := s.repo.GetLatestPredictionByDocumentId(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		prediction = nil
	} else if err != nil {
		return nil, err
	}

	return &DocumentView{
		Document:   document,
		Prediction: prediction,
	}, nil
}
