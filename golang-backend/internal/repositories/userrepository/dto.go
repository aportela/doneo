package userrepository

import (
	"database/sql"

	"github.com/aportela/doneo/internal/repositories"
)

type UserBaseDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type UserDTO struct {
	UserBaseDTO
	Email        string        `db:"email"`
	PasswordHash string        `db:"password_hash"`
	CreatedAt    int64         `db:"created_at"`
	UpdatedAt    sql.NullInt64 `db:"updated_at"`
	DeletedAt    sql.NullInt64 `db:"deleted_at"`
	IsSuperUser  bool          `db:"is_super_user"`
}

type SearchUsersFilterDTO struct {
	Name              *string
	Email             *string
	AdministratorFlag *bool
	CreatedAt         *repositories.TimestampFilter
	UpdatedAt         *repositories.TimestampFilter
	DeletedAt         *repositories.TimestampFilter
}
