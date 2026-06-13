package service

import (
	"io"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
)

type FileStorage interface {
	Save(filename string, src io.Reader) error
}

type UploadResult struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type DocumentService struct {
	fs FileStorage
}

func NewDocumentService(fs FileStorage) *DocumentService {
	return &DocumentService{
		fs: fs,
	}
}

func (s *DocumentService) Upload(
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

	return &UploadResult{
		ID:       documentID,
		Filename: filename,
	}, nil
}
