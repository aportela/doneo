package taskrepository

import (
	"database/sql"

	"github.com/aportela/doneo/internal/repositories"
)

type taskDTO struct {
	ID          string         `db:"id"`
	Index       uint           `db:"task_index"`
	Summary     string         `db:"summary"`
	Description sql.NullString `db:"description"`
	CreatorId   string         `db:"creator_id"`
	CreatorName string         `db:"creator_name"`
	CreatedAt   int64          `db:"created_at"`
	UpdatedAt   sql.NullInt64  `db:"updated_at"`
	DeletedAt   sql.NullInt64  `db:"deleted_at"`
	StartedAt   sql.NullInt64  `db:"started_at"`
	FinishedAt  sql.NullInt64  `db:"finished_at"`
	DueAt       sql.NullInt64  `db:"due_at"`
}

type searchFilterDTO struct {
	Summary         *string
	CreatedAt       *repositories.TimestampFilter
	CreatedByUserId *string
}
