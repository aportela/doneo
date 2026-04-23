package models

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	CreatedAt    int64  `json:"createdAt"`
	LastUpdateAt *int64 `json:"lastUpdateAt"`
}
