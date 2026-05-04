package domain

import "time"

type UserBase struct {
	ID        string
	Name      string
	AvatarURL string
}

type User struct {
	UserBase
	Email        string
	Password     string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
	DeletedAt    *time.Time
	IsSuperUser  bool
}

func (u *User) IsActive() bool {
	return u.DeletedAt == nil
}

func (u *User) IsAdmin() bool {
	return u.IsSuperUser
}
