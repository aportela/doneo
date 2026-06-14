package historyoperationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/utils"
)

type HistoryOperationService interface {
	AddProjectHistoryOperation(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	SearchProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.HistoryOperation, error)
	AddTaskHistoryOperation(ctx context.Context, projectId string, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	SearchTaskHistoryOperations(ctx context.Context, taskId string) ([]domain.HistoryOperation, error)
}

type historyOperationService struct {
	database   database.Database
	repository historyoperationrepository.HistoryOperationRepository
}

func NewService(db database.Database, repository historyoperationrepository.HistoryOperationRepository) HistoryOperationService {
	return &historyOperationService{database: db, repository: repository}
}

func (service *historyOperationService) AddProjectHistoryOperation(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddProjectHistoryOperation(ctx, projectId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) SearchProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.SearchProjectHistoryOperations(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get project history operations: %w", err)
	}
	return operations, nil
}

func (service *historyOperationService) AddTaskHistoryOperation(ctx context.Context, projectId string, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddTaskOperation(ctx, projectId, taskId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) SearchTaskHistoryOperations(ctx context.Context, taskId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.SearchTaskHistoryOperations(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get task history operations: %w", err)
	}
	return operations, nil
}
