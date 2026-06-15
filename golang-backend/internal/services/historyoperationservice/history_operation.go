package historyoperationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/utils"
)

type HistoryOperationService interface {
	AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectId string) ([]domain.HistoryOperation, error)
	AddTaskHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskId string) ([]domain.HistoryOperation, error)
}

type historyOperationService struct {
	repository historyoperationrepository.HistoryOperationRepository
}

func NewService(repository historyoperationrepository.HistoryOperationRepository) HistoryOperationService {
	return &historyOperationService{repository: repository}
}

func (service *historyOperationService) AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddProjectHistoryOperation(ctx, dbExecutor, projectId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.SearchProjectHistoryOperations(ctx, dbExecutor, projectId)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get project history operations: %w", err)
	}
	return operations, nil
}

func (service *historyOperationService) AddTaskHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddTaskOperation(ctx, dbExecutor, projectId, taskId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.SearchTaskHistoryOperations(ctx, dbExecutor, taskId)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get task history operations: %w", err)
	}
	return operations, nil
}
