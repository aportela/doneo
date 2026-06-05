package projecthistoryrepository

type projectHistoryOperationDTO struct {
	UserId        string `db:"user_id"`
	UserName      string `db:"user_name"`
	CreatedAt     int64  `db:"created_at"`
	OperationType uint   `db:"operation_type "`
}
