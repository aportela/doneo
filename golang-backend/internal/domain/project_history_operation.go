package domain

import (
	"time"
)

type ProjectHistoryOperation struct {
	ID            string
	CreatedBy     UserBase
	CreatedAt     time.Time
	OperationType ProjectEventType
}
