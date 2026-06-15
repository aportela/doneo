package historyoperationrepository

type historyOperationDTO struct {
	ID            string `db:"id"`
	UserId        string `db:"operation_user_id"`
	UserName      string `db:"user_name"`
	CreatedAt     int64  `db:"operation_date"`
	OperationType uint   `db:"operation_type "`
}

type scope string

const (
	scopeProject scope = "project"
	scopeTask    scope = "task"
)
