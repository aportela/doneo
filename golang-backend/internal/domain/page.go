package domain

import "time"

type Page struct {
	ID        string
	CreatedBy UserBase
	CreatedAt time.Time
	UpdatedAt *time.Time
	Title     string
	Body      string
}
