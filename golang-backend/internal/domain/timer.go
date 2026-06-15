package domain

import "time"

type UserTimer struct {
	ID         string
	Summary    string
	StartedAt  time.Time
	FinishedAt *time.Time
}
