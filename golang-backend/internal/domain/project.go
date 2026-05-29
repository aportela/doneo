package domain

type Project struct {
	ID                     string
	Key                    string
	Summary                string
	Description            *string
	CreatedBy              UserBase
	CreatedAt              int64
	UpdatedAt              *int64
	StartedAt              *int64
	FinishedAt             *int64
	DueAt                  *int64
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
	CreatedAt       *TimestampFilter
	CreatedByUserId *string
}
