package userrepository

import "database/sql"

type userBaseDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type userDTO struct {
	userBaseDTO
	Email        string        `db:"email"`
	PasswordHash string        `db:"password_hash"`
	CreatedAt    int64         `db:"created_at"`
	UpdatedAt    sql.NullInt64 `db:"updated_at"`
	DeletedAt    sql.NullInt64 `db:"deleted_at"`
	IsSuperUser  bool          `db:"is_super_user"`
}
