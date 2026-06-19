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

type ProjectEventType uint

const (
	EventProjectCreated ProjectEventType = 100
	EventProjectUpdated ProjectEventType = 101
	EventProjectDeleted ProjectEventType = 102

	EventProjectNoteAdded   ProjectEventType = 110
	EventProjectNoteUpdated ProjectEventType = 111
	EventProjectNoteDeleted ProjectEventType = 112

	EventProjectAttachmentAdded   ProjectEventType = 120
	EventProjectAttachmentDeleted ProjectEventType = 122

	EventProjectPermissionAdded   ProjectEventType = 130
	EventProjectPermissionDeleted ProjectEventType = 132

	EventTaskCreated ProjectEventType = 200
	EventTaskUpdated ProjectEventType = 201
	EventTaskDeleted ProjectEventType = 202

	EventTaskNoteAdded   ProjectEventType = 210
	EventTaskNoteUpdated ProjectEventType = 211
	EventTaskNoteDeleted ProjectEventType = 212

	EventTaskAttachmentAdded   ProjectEventType = 220
	EventTaskAttachmentDeleted ProjectEventType = 222

	EventTaskTimeEntryAdded   ProjectEventType = 230
	EventTaskTimeEntryUpdated ProjectEventType = 231
	EventTaskTimeEntryDeleted ProjectEventType = 232
)
