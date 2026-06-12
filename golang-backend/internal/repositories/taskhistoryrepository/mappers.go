package taskhistoryrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
)

func toDTO(operation domain.HistoryOperation) taskHistoryOperationDTO {
	return taskHistoryOperationDTO{
		ID:            operation.ID,
		UserId:        operation.CreatedBy.ID,
		UserName:      operation.CreatedBy.Name,
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func toDomain(operation taskHistoryOperationDTO) domain.HistoryOperation {
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

func toDomainArray(operations []taskHistoryOperationDTO) []domain.HistoryOperation {
	results := make([]domain.HistoryOperation, 0, len(operations))
	for _, operation := range operations {
		results = append(results, toDomain(operation))
	}
	return results
}
