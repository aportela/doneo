package pagerepository

import "database/sql"

type pageDTO struct {
	ID          string        `db:"id"`
	CreatorID   string        `db:"creator_id"`
	CreatorName string        `db:"creator_name"`
	CreatedAt   int64         `db:"created_at"`
	UpdatedAt   sql.NullInt64 `db:"updated_at"`
	Title       string        `db:"title"`
	Body        string        `db:"body"`
}
