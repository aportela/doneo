package domain

import "time"

type Project struct {
	ID                     string
	Key                    string
	Summary                string
	Description            *string
	CreatedBy              UserBase
	CreatedAt              time.Time
	UpdatedAt              *time.Time
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
	Key             *string
	Summary         *string
	TypeId          *string
	StatusId        *string
	PriorityId      *string
	CreatedAt       *TimestampFilter
	CreatedByUserId *string
}

type ProjectEventType uint

const (
	EventProjectCreated ProjectEventType = 1
	EventProjectUpdated ProjectEventType = 2
	EventProjectDeleted ProjectEventType = 3
)
