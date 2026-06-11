package domain

import "time"

type Timer struct {
	ID         string
	Summary    string
	StartedAt  time.Time
	FinishedAt *time.Time
}
