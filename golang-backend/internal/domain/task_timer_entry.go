package domain

import "time"

type TaskTimerEntry struct {
	ID           string
	Summary      string
	CreatedBy    UserBase
	CreatedAt    time.Time
	TotalSeconds uint
}
