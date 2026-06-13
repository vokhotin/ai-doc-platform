package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalFileStorage struct {
	uploadsDir string
}

func NewLocalFileStorage(uploadsDir string) *LocalFileStorage {
	return &LocalFileStorage{uploadsDir: uploadsDir}
}

func (s *LocalFileStorage) Save(filename string, src io.Reader) error {
	dst, err := os.Create(filepath.Join(s.uploadsDir, filename))

	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)
	if err != nil {
		dst.Close()
		return err
	}

	return dst.Close()
}
