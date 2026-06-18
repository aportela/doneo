package domain

import (
	"time"
)

type Attachment struct {
	ID           string
	CreatedBy    UserBase
	CreatedAt    time.Time
	OriginalName string
	ContentType  string
	Size         uint32
}
