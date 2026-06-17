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

type mockDocumentRepository struct {
	saveErr error
}

type mockInferenceService struct {
	documentID string
	text       string
	predictErr error
}

type mockCloserReader struct {
	io.ReadCloser
}

func (m *mockDocumentRepository) SaveDocument(ctx context.Context, doc *model.Document) error {
	return m.saveErr
}

func (m *mockDocumentRepository) GetByID(ctx context.Context, id string) (*model.Document, error) {
	return nil, nil
}

func (m *mockInferenceService) Predict(ctx context.Context, documentID string, text string) (*model.InferenceResult, error) {
	m.documentID = documentID
	m.text = text
	return nil, m.predictErr
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
		repo      *mockDocumentRepository
		inference *mockInferenceService
		wantErr   bool
	}{
		{"success", "test.pdf", &mockFileStorage{}, &mockDocumentRepository{}, &mockInferenceService{}, false},
		{"success empty file extension", "test", &mockFileStorage{}, &mockDocumentRepository{}, &mockInferenceService{}, false},
		{"storage failure", "test.pdf", &mockFileStorage{saveErr: errors.New("disk is full")}, &mockDocumentRepository{}, &mockInferenceService{}, true},
		{"repo failure", "test.pdf", &mockFileStorage{}, &mockDocumentRepository{saveErr: errors.New("db is down")}, &mockInferenceService{}, true},
		{"inference service failure", "test.pdf", &mockFileStorage{}, &mockDocumentRepository{}, &mockInferenceService{predictErr: errors.New("inference service is down")}, true},
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

			if result.ID != tt.inference.documentID {
				t.Errorf("expected documentID %s, got %s", result.ID, tt.inference.documentID)
			}

			//TODO remove filename as text and test real extraction.
		})
	}

}
