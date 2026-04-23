package models

type UserBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	UserBase
	Email        string  `json:"email"`
	Password     *string `json:"string"`
	CreatedAt    int64   `json:"createdAt"`
	LastUpdateAt *int64  `json:"lastUpdateAt"`
}
