package domain

type Attachment struct {
	ID        string   `json:"id"`
	FileName  string   `json:"fileName"`
	SizeBytes int64    `json:"sizeBytes"`
	CreatedBy UserBase `json:"createdBy"`
	CreatedAt int64    `json:"createdAt"`
}
