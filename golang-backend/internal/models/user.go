package models

type UserBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	UserBase
	Email       string  `json:"email"`
	Password    *string `json:"password"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   *int64  `json:"updatedAt"`
	IsSuperUser bool    `json:"isSuperUser"`
	AvatarURL   string  `json:"avatar"`
}
