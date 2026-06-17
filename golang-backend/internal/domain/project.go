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
	EventProjectCreated           ProjectEventType = 1
	EventProjectUpdated           ProjectEventType = 2
	EventProjectDeleted           ProjectEventType = 3
	EventProjectNoteAdded         ProjectEventType = 4
	EventProjectNoteUpdated       ProjectEventType = 5
	EventProjectNoteDeleted       ProjectEventType = 6
	EventProjectAttachmentAdded   ProjectEventType = 7
	EventProjectAttachmentDeleted ProjectEventType = 8
	EventProjectPermissionAdded   ProjectEventType = 9
	EventProjectPermissionDeleted ProjectEventType = 10

	EventTaskCreated           ProjectEventType = 20
	EventTaskUpdated           ProjectEventType = 21
	EventTaskDeleted           ProjectEventType = 22
	EventTaskNoteAdded         ProjectEventType = 23
	EventTaskNoteUpdated       ProjectEventType = 24
	EventTaskNoteDeleted       ProjectEventType = 25
	EventTaskAttachmentAdded   ProjectEventType = 26
	EventTaskAttachmentDeleted ProjectEventType = 27
	EventTaskTimeEntryAdded    ProjectEventType = 28
	EventTaskTimeEntryUpdated  ProjectEventType = 29
	EventTaskTimeEntryDeleted  ProjectEventType = 30
)
