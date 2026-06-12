package taskhistoryhandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type TaskHistoryOperationResponse struct {
	ID            string                       `json:"id"`
	CreatedBy     userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt     int64                        `json:"createdAt"`
	OperationType uint                         `json:"operationType"`
}

type searchResponse struct {
	HistoryOperations []TaskHistoryOperationResponse `json:"historyOperations"`
}
