package domain

import "time"

type Project struct {
	ID                     string
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
	Type                   ProjectType
	Priority               ProjectPriority
	Status                 ProjectStatus
	TasksCount             uint
	PermissionsCount       uint
	AttachmentsCount       uint
	NotesCount             uint
	HistoryOperationsCount uint

	//lead, asignee
}

type SearchProjectFilter struct {
	Slug            *string
	Summary         *string
	TypeID          *string
	StatusID        *string
	PriorityID      *string
	CreatedAt       *TimestampFilter
	CreatedByUserID *string
	ViewByUserID    *string
}
