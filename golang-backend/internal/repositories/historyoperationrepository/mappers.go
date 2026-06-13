package historyoperationrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
)

func toDTO(operation domain.HistoryOperation) historyOperationDTO {
	return historyOperationDTO{
		ID:            operation.ID,
		UserId:        operation.CreatedBy.ID,
		UserName:      operation.CreatedBy.Name,
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func toDomain(operation historyOperationDTO) domain.HistoryOperation {
	switch operation.OperationType {
	case 1:
		break
	}
	return domain.HistoryOperation{
		ID: operation.ID,
		CreatedBy: domain.UserBase{
			ID:   operation.UserId,
			Name: operation.UserName,
		},
		CreatedAt:     time.UnixMilli(operation.CreatedAt),
		OperationType: domain.ProjectEventType(operation.OperationType),
	}
}

func toDomainArray(operations []historyOperationDTO) []domain.HistoryOperation {
	results := make([]domain.HistoryOperation, 0, len(operations))
	for _, operation := range operations {
		results = append(results, toDomain(operation))
	}
	return results
}
