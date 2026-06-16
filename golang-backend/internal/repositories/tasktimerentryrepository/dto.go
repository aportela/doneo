package tasktimerentryrepository

type taskTimerEntryDTO struct {
	ID           string `db:"id"`
	CreatorId    string `db:"creator_id"`
	CreatorName  string `db:"creator_name"`
	CreatedAt    int64  `db:"created_at"`
	Summary      string `db:"summary"`
	TotalSeconds uint   `db:"total_seconds"`
}
