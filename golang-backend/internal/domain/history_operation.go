package domain

import (
	"time"
)

type HistoryOperationEventType uint16

const (
	EventProjectCreated HistoryOperationEventType = 100
	EventProjectUpdated HistoryOperationEventType = 101
	EventProjectDeleted HistoryOperationEventType = 102

	EventProjectNoteAdded   HistoryOperationEventType = 110
	EventProjectNoteUpdated HistoryOperationEventType = 111
	EventProjectNoteDeleted HistoryOperationEventType = 112

	EventProjectAttachmentAdded   HistoryOperationEventType = 120
	EventProjectAttachmentDeleted HistoryOperationEventType = 122

	EventProjectPermissionAdded   HistoryOperationEventType = 130
	EventProjectPermissionDeleted HistoryOperationEventType = 132

	EventTaskCreated HistoryOperationEventType = 200
	EventTaskUpdated HistoryOperationEventType = 201
	EventTaskDeleted HistoryOperationEventType = 202

	EventTaskNoteAdded   HistoryOperationEventType = 210
	EventTaskNoteUpdated HistoryOperationEventType = 211
	EventTaskNoteDeleted HistoryOperationEventType = 212

	EventTaskAttachmentAdded   HistoryOperationEventType = 220
	EventTaskAttachmentDeleted HistoryOperationEventType = 222

	EventTaskTimeEntryAdded   HistoryOperationEventType = 230
	EventTaskTimeEntryUpdated HistoryOperationEventType = 231
	EventTaskTimeEntryDeleted HistoryOperationEventType = 232
)

type HistoryOperation struct {
	ID            string
	CreatedBy     UserBase
	CreatedAt     time.Time
	OperationType HistoryOperationEventType
}
