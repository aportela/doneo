package models

type Note struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	CreatedBy User   `json:"createdBy"`
	CreatedAt int64  `json:"createdAt"`
}
