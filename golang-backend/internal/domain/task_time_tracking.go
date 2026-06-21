package domain

import "time"

type TaskTimeTracking struct {
	ID        string
	Summary   string
	CreatedBy UserBase
	CreatedAt time.Time
	SpentTime uint64
}
