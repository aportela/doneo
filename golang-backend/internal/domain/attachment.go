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

const MaxUploadSize = 32 << 20 // 32 MiB
