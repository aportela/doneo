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
	AddProjectHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	GetProjectHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.HistoryOperation, error)
	AddTaskHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	GetTaskHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.HistoryOperation, error)
}

type historyOperationService struct {
	repository historyoperationrepository.HistoryOperationRepository
}

func NewService(repository historyoperationrepository.HistoryOperationRepository) HistoryOperationService {
	return &historyOperationService{repository: repository}
}

func (service *historyOperationService) AddProjectHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddProjectHistoryOperation(ctx, dbExecutor, projectID, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) GetProjectHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.GetProjectHistoryOperations(ctx, dbExecutor, projectID)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get project history operations: %w", err)
	}
	return operations, nil
}

func (service *historyOperationService) AddTaskHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddTaskHistoryOperation(ctx, dbExecutor, projectID, taskID, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) GetTaskHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.GetTaskHistoryOperations(ctx, dbExecutor, taskID)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get task history operations: %w", err)
	}
	return operations, nil
}
