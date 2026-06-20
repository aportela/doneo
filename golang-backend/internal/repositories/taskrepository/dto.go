package taskrepository

import (
	"database/sql"

	"github.com/aportela/doneo/internal/repositories"
)

type taskDTO struct {
	ID                     string         `db:"id"`
	projectID              string         `db:"project_id"`
	Index                  uint           `db:"task_index"`
	Slug                   string         `db:"task_slug"` // TODO: projectSlug (without Index ???)
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
	StatusID               string         `db:"status_id"`
	StatusName             string         `db:"status_name"`
	StatusHexColor         string         `db:"status_hex_color"`
	PriorityID             string         `db:"priority_id"`
	PriorityName           string         `db:"priority_name"`
	PriorityHexColor       string         `db:"priority_hex_color"`
	AttachmentsCount       uint16         `db:"attachments_count"`
	NotesCount             uint16         `db:"notes_count"`
	HistoryOperationsCount uint16         `db:"history_operations_count"`
	TimeTrackingsCount     uint16         `db:"time_trackings_count"`
}

type searchFilterDTO struct {
	ProjectID       *string // TODO ?????
	Summary         *string
	PriorityID      *string
	StatusID        *string
	CreatedAt       *repositories.TimestampFilter
	CreatedByUserID *string
	ViewByUserID    *string
}
