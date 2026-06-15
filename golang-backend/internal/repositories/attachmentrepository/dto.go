package attachmentrepository

type attachmentDTO struct {
	ID           string `db:"id"`
	CreatorID    string `db:"creator_id"`
	UserName     string `db:"user_name"`
	CreatedAt    int64  `db:"created_at"`
	OriginalName string `db:"original_name"`
	ContentType  string `db:"content_type"`
	Size         uint32 `db:"size"`
}
