package service

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type mockFileStorage struct {
	saveErr       error
	savedFilename string
}

func (m *mockFileStorage) Save(filename string, src io.Reader) error {
	m.savedFilename = filename
	return m.saveErr
}

type mockRepository struct {
	prediction     *model.Prediction
	saveErr        error
	updateDocError error
	txError        error
}

func (m *mockRepository) UpdateDocumentStatus(ctx context.Context, documentID string, status model.DocumentStatus) error {
	return m.updateDocError
}

func (m *mockRepository) SavePredictionAndMarkDocumentDone(ctx context.Context, prediction *model.Prediction) error {
	return m.txError
}

func (m *mockRepository) GetLatestPredictionByDocumentId(ctx context.Context, documentId string) (*model.Prediction, error) {
	return m.prediction, m.saveErr
}

type mockInferenceService struct {
	*model.InferenceResult
	predictErr error
}

type mockCloserReader struct {
	io.ReadCloser
}

func (m *mockRepository) SaveDocument(ctx context.Context, doc *model.Document) error {
	return m.saveErr
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	return nil, nil
}

func (m *mockInferenceService) Predict(ctx context.Context, documentID string, text string) (*model.InferenceResult, error) {
	m.InferenceResult = &model.InferenceResult{
		DocumentID: documentID,
		Label:      "finance",
		Confidence: 0.8,
	}
	return m.InferenceResult, m.predictErr
}

func (rc mockCloserReader) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

func (rc mockCloserReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func TestDocumentService_Upload(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		storage   *mockFileStorage
		repo      *mockRepository
		inference *mockInferenceService
		wantErr   bool
	}{
		{"success", "test.pdf", &mockFileStorage{}, &mockRepository{}, &mockInferenceService{}, false},
		{"success empty file extension", "test", &mockFileStorage{}, &mockRepository{}, &mockInferenceService{}, false},
		{"storage failure", "test.pdf", &mockFileStorage{saveErr: errors.New("disk is full")}, &mockRepository{}, &mockInferenceService{}, true},
		{"repo failure", "test.pdf", &mockFileStorage{}, &mockRepository{saveErr: errors.New("db is down")}, &mockInferenceService{}, true},
		{"inference service failure", "test.pdf", &mockFileStorage{}, &mockRepository{}, &mockInferenceService{predictErr: errors.New("inference service is down")}, true},
		{"inference service failure and failed to update document status", "test.pdf", &mockFileStorage{}, &mockRepository{updateDocError: errors.New("failed update error")}, &mockInferenceService{predictErr: errors.New("inference service is down")}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			svc := NewDocumentService(tt.storage, tt.repo, tt.inference)
			mc := &mockCloserReader{io.NopCloser(strings.NewReader("data"))}
			result, err := svc.Upload(ctx, mc, tt.filename)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expect an error got nothing")
				}
				if !(errors.Is(err, tt.storage.saveErr) || errors.Is(err, tt.repo.saveErr) || errors.Is(err, tt.inference.predictErr)) {
					t.Fatalf("expect specific error got %v", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("expect no error, got %v", err)
			}

			if result.Filename != tt.filename {
				t.Fatalf("expect %v got %v", tt.filename, result.Filename)
			}

			if result.ID == "" {
				t.Error("expected non-empty ID")
			}

			if filepath.Ext(tt.storage.savedFilename) != filepath.Ext(tt.filename) {
				t.Errorf("expected extension %s, got %s", filepath.Ext(tt.filename),
					filepath.Ext(tt.storage.savedFilename))
			}

			if result.ID != tt.inference.InferenceResult.DocumentID {
				t.Errorf("expected documentID %s, got %s", result.ID, tt.inference.InferenceResult.DocumentID)
			}

			//TODO remove filename as text and test real extraction.
		})
	}

}
