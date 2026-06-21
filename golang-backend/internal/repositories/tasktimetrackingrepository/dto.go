package tasktimetrackingrepository

type taskTimeTrackingDTO struct {
	ID          string `db:"id"`
	CreatorID   string `db:"creator_id"`
	CreatorName string `db:"creator_name"`
	CreatedAt   int64  `db:"created_at"`
	Summary     string `db:"summary"`
	SpentTime   uint64 `db:"spent_time"`
}
