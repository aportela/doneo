package historyoperationhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func domainToResponse(operation domain.HistoryOperation) HistoryOperationResponse {
	return HistoryOperationResponse{
		ID:            operation.ID,
		CreatedBy:     userhandler.BaseDomainToBaseResponse(operation.CreatedBy),
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint16(operation.OperationType),
	}
}

func domainArrayToResponseArray(attachments []domain.HistoryOperation) []HistoryOperationResponse {
	operationResponse := []HistoryOperationResponse{}
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
