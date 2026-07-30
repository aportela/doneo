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
	ArchivedAt             *time.Time
	Type                   ProjectType
	Priority               ProjectPriority
	Status                 ProjectStatus
	TasksCount             uint16
	PermissionsCount       uint16
	AttachmentsCount       uint16
	NotesCount             uint16
	PagesCount             uint16
	HistoryOperationsCount uint16
	PermissionsBitMask     Bitmask
	//lead, asignee
}

type SearchProjectFilter struct {
	Slug            *string
	Summary         *string
	TypeID          *string
	StatusID        *string
	PriorityID      *string
	ViewByUserID    *string
	CreatedByUserID *string
	CreatedAt       *TimestampFilter
	UpdatedAt       *TimestampFilter
	StartedAt       *TimestampFilter
	FinishedAt      *TimestampFilter
	DueAt           *TimestampFilter
	ArchivedAt      *TimestampFilter
	DeletedAt       *TimestampFilter
}
