package models

type Note struct {
	ID        string `json:"id"`
	Body      string `json:"body"`
	CreatedBy User   `json:"createdBy"`
	CreatedAt int64  `json:"createdAt"`
}
