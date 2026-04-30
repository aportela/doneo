package domain

type Project struct {
	ID          string
	WorkspaceId string
	Key         string
	Summary     string
	Description *string
	CreatedBy   UserBase
	CreatedAt   int64
	UpdatedAt   *int64
	StartedAt   *int64
	FinishedAt  *int64
	DueAt       *int64
	Type        ProjectType
	Priority    ProjectPriority
	Status      ProjectStatus
	//lead, asignee
}
