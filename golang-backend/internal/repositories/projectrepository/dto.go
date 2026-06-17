package projectrepository

import (
	"database/sql"

	"github.com/aportela/doneo/internal/repositories"
)

type projectDTO struct {
	ID                     string         `db:"id"`
	Slug                   string         `db:"slug"`
	Summary                string         `db:"summary"`
	Description            sql.NullString `db:"description"`
	CreatorID              string         `db:"creator_id"`
	CreatorName            string         `db:"creator_name"`
	CreatedAt              int64          `db:"created_at"`
	UpdatedAt              sql.NullInt64  `db:"updated_at"`
	DeletedAt              sql.NullInt64  `db:"deleted_at"`
	StartedAt              sql.NullInt64  `db:"started_at"`
	FinishedAt             sql.NullInt64  `db:"finished_at"`
	DueAt                  sql.NullInt64  `db:"due_at"`
	TypeID                 string         `db:"type_id"`
	TypeName               string         `db:"type_name"`
	TypeHexColor           string         `db:"type_hex_color"`
	StatusID               string         `db:"status_id"`
	StatusName             string         `db:"status_name"`
	StatusHexColor         string         `db:"status_hex_color"`
	PriorityID             string         `db:"priority_id"`
	PriorityName           string         `db:"priority_name"`
	PriorityHexColor       string         `db:"priority_hex_color"`
	TasksCount             uint           `db:"tasks_count"`
	PermissionsCount       uint           `db:"permissions_count"`
	AttachmentsCount       uint           `db:"attachments_count"`
	NotesCount             uint           `db:"notes_count"`
	HistoryOperationsCount uint           `db:"history_operations_count"`
}

type searchFilterDTO struct {
	Slug            *string
	Summary         *string
	TypeID          *string
	PriorityID      *string
	StatusID        *string
	CreatedAt       *repositories.TimestampFilter
	CreatedByUserID *string
	ViewByUserID    *string
}
