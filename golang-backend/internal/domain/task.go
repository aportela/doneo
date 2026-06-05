package domain

import "time"

type Task struct {
	ID                     string
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
	PermissionsCount       uint
	AttachmentsCount       uint
	NotesCount             uint
	HistoryOperationsCount uint
}

type SearchTaskFilter struct {
	Summary         *string
	StatusId        *string
	PriorityId      *string
	CreatedAt       *TimestampFilter
	CreatedByUserId *string
}
