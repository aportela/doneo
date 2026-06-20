package domain

import (
	"time"
)

type HistoryOperation struct {
	ID            string
	CreatedBy     UserBase
	CreatedAt     time.Time
	OperationType ProjectEventType
}
