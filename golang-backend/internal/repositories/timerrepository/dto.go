package timerrepository

import "database/sql"

type timerDTO struct {
	ID         string        `db:"id"`
	StartedAt  int64         `db:"started_at"`
	FinishedAt sql.NullInt64 `db:"finished_at"`
}
