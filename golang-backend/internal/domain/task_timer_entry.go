package domain

import "time"

type TaskTimeEntry struct {
	ID           string
	Summary      string
	CreatedBy    UserBase
	CreatedAt    time.Time
	TotalSeconds uint
}
