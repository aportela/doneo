package domain

import "time"

type Task struct {
	ID                     string
	ProjectID              string
	Index                  uint
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
	Priority               TaskPriority
	Status                 TaskStatus
	Tags                   []string
	PermissionsCount       uint
	AttachmentsCount       uint
	NotesCount             uint
	HistoryOperationsCount uint
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
