package projecthistoryhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func domainToResponse(operation domain.ProjectHistoryOperation) ProjectHistoryOperationResponse {
	return ProjectHistoryOperationResponse{
		CreatedBy:     userhandler.BaseDomainToBaseResponse(operation.CreatedBy),
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func domainArrayToResponseArray(attachments []domain.ProjectHistoryOperation) []ProjectHistoryOperationResponse {
	operationResponse := []ProjectHistoryOperationResponse{}
	for _, attachment := range attachments {
		operationResponse = append(operationResponse, domainToResponse(attachment))
	}
	return operationResponse
}

func toSearchResponse(operations []domain.ProjectHistoryOperation) searchResponse {
	return searchResponse{
		HistoryOperations: domainArrayToResponseArray(operations),
	}
}
