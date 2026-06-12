package domain

import (
	"time"
)

// TODO: rename file
type HistoryOperation struct {
	ID            string
	CreatedBy     UserBase
	CreatedAt     time.Time
	OperationType ProjectEventType
}
