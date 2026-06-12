package taskhistoryhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func domainToResponse(operation domain.HistoryOperation) TaskHistoryOperationResponse {
	return TaskHistoryOperationResponse{
		ID:            operation.ID,
		CreatedBy:     userhandler.BaseDomainToBaseResponse(operation.CreatedBy),
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func domainArrayToResponseArray(attachments []domain.HistoryOperation) []TaskHistoryOperationResponse {
	operationResponse := []TaskHistoryOperationResponse{}
	for _, attachment := range attachments {
		operationResponse = append(operationResponse, domainToResponse(attachment))
	}
	return operationResponse
}

func toSearchResponse(operations []domain.HistoryOperation) searchResponse {
	return searchResponse{
		HistoryOperations: domainArrayToResponseArray(operations),
	}
}
