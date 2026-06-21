package domain

import "time"

type Task struct {
	ID                     string
	ProjectID              string
	Index                  uint16
	Slug                   string
	Summary                string
	Description            *string
	CreatedBy              UserBase
	CreatedAt              time.Time
	UpdatedAt              *time.Time
	DeletedAt              *time.Time
	StartedAt              *time.Time
	FinishedAt             *time.Time
	DueAt                  *time.Time
	EstimatedTime          uint64
	TotalSpentTime         uint64
	Priority               TaskPriority
	Status                 TaskStatus
	Tags                   []string
	PermissionsCount       uint16
	AttachmentsCount       uint16
	NotesCount             uint16
	HistoryOperationsCount uint16
	TimeTrackingsCount     uint16
}

type SearchTaskFilter struct {
	ProjectID       *string // TODO ???
	Summary         *string
	StatusID        *string
	PriorityID      *string
	CreatedAt       *TimestampFilter
	CreatedByUserID *string
	ViewByUserID    *string
}
