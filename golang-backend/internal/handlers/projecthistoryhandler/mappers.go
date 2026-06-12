package projecthistoryhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func domainToResponse(operation domain.HistoryOperation) ProjectHistoryOperationResponse {
	return ProjectHistoryOperationResponse{
		ID:            operation.ID,
		CreatedBy:     userhandler.BaseDomainToBaseResponse(operation.CreatedBy),
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func domainArrayToResponseArray(attachments []domain.HistoryOperation) []ProjectHistoryOperationResponse {
	operationResponse := []ProjectHistoryOperationResponse{}
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
