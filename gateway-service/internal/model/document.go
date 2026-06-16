package model

import (
	"time"
)

type Document struct {
	ID               string
	OriginalFilename string
	StoredFilename   string
	Status           DocumentStatus
	CreatedAt        time.Time
}

type DocumentStatus string

const (
	StatusPending    DocumentStatus = "pending"
	StatusProcessing DocumentStatus = "processing"
	StatusDone       DocumentStatus = "done"
)
