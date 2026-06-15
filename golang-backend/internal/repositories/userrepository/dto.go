package userrepository

import (
	"database/sql"

	"github.com/aportela/doneo/internal/repositories"
)

type userBaseDTO struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type userDTO struct {
	userBaseDTO
	Email              string        `db:"email"`
	PasswordHash       string        `db:"password_hash"`
	CreatedAt          int64         `db:"created_at"`
	UpdatedAt          sql.NullInt64 `db:"updated_at"`
	DeletedAt          sql.NullInt64 `db:"deleted_at"`
	PermissionsBitmask uint64        `db:"permissions_bitmask"`
}

type searchFilterDTO struct {
	Name                        *string
	Email                       *string
	RequiredPermissionsBitmask  *uint64
	ForbiddenPermissionsBitmask *uint64
	CreatedAt                   *repositories.TimestampFilter
	UpdatedAt                   *repositories.TimestampFilter
	DeletedAt                   *repositories.TimestampFilter
}
