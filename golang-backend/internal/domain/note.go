package domain

import "time"

type Note struct {
	ID        string
	CreatedBy UserBase
	CreatedAt time.Time
	UpdatedAt *time.Time
	Body      string
}
