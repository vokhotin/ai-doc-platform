package model

import "time"

type Prediction struct {
	Id         string
	DocumentID string
	Label      string
	Confidence float32
	CreatedAt  time.Time
}
