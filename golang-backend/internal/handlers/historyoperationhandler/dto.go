package historyoperationhandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type historyOperationResponse struct {
	ID            string                       `json:"id"`
	CreatedBy     userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt     int64                        `json:"createdAt"`
	OperationType uint16                       `json:"operationType"`
}

type searchResponse struct {
	HistoryOperations []historyOperationResponse `json:"historyOperations"`
}
