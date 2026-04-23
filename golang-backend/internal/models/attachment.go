package models

type Attachment struct {
	ID        int    `json:"id"`
	FileName  string `json:"fileName"`
	SizeBytes int64  `json:"sizeBytes"`
	CreatedAt int64  `json:"createdAt"`
}
