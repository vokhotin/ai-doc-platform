package service

import (
	"context"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/vokhotin/ai-doc-platform/gateway-service/internal/model"
)

type FileStorage interface {
	Save(filename string, src io.Reader) error
}

type DocumentRepository interface {
	Save(ctx context.Context, doc *model.Document) error
}

type UploadResult struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type DocumentService struct {
	fs FileStorage
	dr DocumentRepository
}

func NewDocumentService(fs FileStorage, dr DocumentRepository) *DocumentService {
	return &DocumentService{
		fs: fs,
		dr: dr,
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
	err = s.dr.Save(ctx, doc)
	if err != nil {
		return nil, err
	}

	return &UploadResult{
		ID:       documentID,
		Filename: filename,
	}, nil
}
