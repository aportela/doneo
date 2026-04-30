package projectrepository

type projectDTO struct {
	ID           string  `db:"id"`
	WorkspaceId  string  `db:"workspace_id"`
	Key          string  `db:"key"`
	Summary      string  `db:"summary"`
	Description  *string `json:"description"`
	Status       uint8   `json:"status"`
	Priority     uint8   `json:"priority"`
	CreatorId    string  `db:"creator_id"`
	CreatorName  string  `db:"creator_name"`
	CreatedAt    int64   `db:"created_at"`
	UpdatedAt    *int64  `db:"updated_at"`
	StartedAt    *int64  `db:"started_at"`
	FinishedAt   *int64  `db:"finished_at"`
	DueAt        *int64  `db:"due_at"`
	TypeId       string  `db:"type_id"`
	TypeName     string  `db:"type_name"`
	StatusId     string  `db:"status_id"`
	StatusName   string  `db:"status_name"`
	PriorityId   string  `db:"priority_id"`
	PriorityName string  `db:"priority_name"`
}
