package domain

import (
	"time"
)

type ProjectHistoryOperation struct {
	CreatedBy     UserBase
	CreatedAt     time.Time
	OperationType ProjectEventType
}
