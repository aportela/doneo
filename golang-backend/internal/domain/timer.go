package domain

import "time"

type Timer struct {
	ID         string
	StartedAt  time.Time
	FinishedAt *time.Time
}
