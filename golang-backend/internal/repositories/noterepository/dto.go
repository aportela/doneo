package noterepository

import "database/sql"

type noteDTO struct {
	ID          string        `db:"id"`
	CreatorID   string        `db:"creator_id"`
	CreatorName string        `db:"creator_name"`
	CreatedAt   int64         `db:"created_at"`
	UpdatedAt   sql.NullInt64 `db:"updated_at"`
	Body        string        `db:"body"`
}
