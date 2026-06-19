package model

import "time"

type Prediction struct {
	ID         string
	DocumentID string
	Label      string
	Confidence float64
	CreatedAt  time.Time
}
