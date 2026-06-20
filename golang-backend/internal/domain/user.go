package domain

import (
	"time"
)

const (
	UserPermissionAdmin Bitmask = 1 << iota
)

type UserBase struct {
	ID   string
	Name string
}

type User struct {
	UserBase
	Email              string
	PasswordHash       string
	CreatedAt          time.Time
	UpdatedAt          *time.Time
	DeletedAt          *time.Time
	PermissionsBitmask Bitmask
}

func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}

func (u *User) IsAdmin() bool {
	return u.PermissionsBitmask.HasFlag(UserPermissionAdmin)
}

type SearchUsersFilter struct {
	Name                        *string
	Email                       *string
	RequiredPermissionsBitmask  *Bitmask
	ForbiddenPermissionsBitmask *Bitmask
	CreatedAt                   *TimestampFilter
	UpdatedAt                   *TimestampFilter
	DeletedAt                   *TimestampFilter
}
