package projecthistoryrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
)

func toDTO(operation domain.ProjectHistoryOperation) projectHistoryOperationDTO {
	return projectHistoryOperationDTO{
		ID:            operation.ID,
		UserId:        operation.CreatedBy.ID,
		UserName:      operation.CreatedBy.Name,
		CreatedAt:     operation.CreatedAt.UnixMilli(),
		OperationType: uint(operation.OperationType),
	}
}

func toDomain(operation projectHistoryOperationDTO) domain.ProjectHistoryOperation {
	switch operation.OperationType {
	case 1:
		break
	}
	return domain.ProjectHistoryOperation{
		ID: operation.ID,
		CreatedBy: domain.UserBase{
			ID:   operation.UserId,
			Name: operation.UserName,
		},
		CreatedAt:     time.UnixMilli(operation.CreatedAt),
		OperationType: domain.ProjectEventType(operation.OperationType),
	}
}

func toDomainArray(operations []projectHistoryOperationDTO) []domain.ProjectHistoryOperation {
	results := make([]domain.ProjectHistoryOperation, 0, len(operations))
	for _, operation := range operations {
		results = append(results, toDomain(operation))
	}
	return results
}
