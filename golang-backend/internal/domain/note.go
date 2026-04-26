package domain

type Note struct {
	ID        string   `json:"id"`
	Body      string   `json:"body"`
	CreatedBy UserBase `json:"createdBy"`
	CreatedAt int64    `json:"createdAt"`
}
