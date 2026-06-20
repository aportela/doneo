package historyoperationrepository

type historyOperationDTO struct {
	ID            string `db:"id"`
	UserID        string `db:"operation_user_id"`
	UserName      string `db:"user_name"`
	CreatedAt     int64  `db:"operation_date"`
	OperationType uint16 `db:"operation_type "`
}

type scope string

const (
	scopeProject scope = "project"
	scopeTask    scope = "task"
)
