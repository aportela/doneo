package usertimerrepository

import "database/sql"

type userTimerDTO struct {
	ID         string        `db:"id"`
	Summary    string        `db:"summary"`
	StartedAt  int64         `db:"started_at"`
	FinishedAt sql.NullInt64 `db:"finished_at"`
}
