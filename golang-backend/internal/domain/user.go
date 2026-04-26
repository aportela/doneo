package domain

type UserBase struct {
	ID   string
	Name string
}

type User struct {
	UserBase
	Email       string
	Password    *string
	CreatedAt   int64
	UpdatedAt   *int64
	IsSuperUser bool
	AvatarURL   string
}
