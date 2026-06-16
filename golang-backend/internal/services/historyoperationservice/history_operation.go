package historyoperationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type HistoryOperationService interface {
	AddProjectHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	GetProjectHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.HistoryOperation, error)
	AddTaskHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	GetTaskHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string) ([]domain.HistoryOperation, error)
}

type historyOperationService struct {
	authorizationService       authorizationservice.AuthorizationService
	historyOperationRepository historyoperationrepository.HistoryOperationRepository
}

func NewService(authorizationService authorizationservice.AuthorizationService, historyOperationRepository historyoperationrepository.HistoryOperationRepository) HistoryOperationService {
	return &historyOperationService{authorizationService: authorizationService, historyOperationRepository: historyOperationRepository}
}

func (service *historyOperationService) AddProjectHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.historyOperationRepository.AddProjectHistoryOperation(ctx, dbExecutor, projectID, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) GetProjectHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.HistoryOperation, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return nil, fmt.Errorf("[HistoryOperationService] user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, contextUser.ID, projectID); err != nil {
		return nil, err
	}
	operations, err := service.historyOperationRepository.GetProjectHistoryOperations(ctx, dbExecutor, projectID)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get project history operations: %w", err)
	}
	return operations, nil
}

func (service *historyOperationService) AddTaskHistoryOperation(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.historyOperationRepository.AddTaskHistoryOperation(ctx, dbExecutor, projectID, taskID, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) GetTaskHistoryOperations(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, taskID string) ([]domain.HistoryOperation, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return nil, fmt.Errorf("[HistoryOperationService] user not found in context")
	}
	if err := service.authorizationService.RequireTaskViewPermission(ctx, contextUser.ID, projectID); err != nil {
		return nil, err
	}
	operations, err := service.historyOperationRepository.GetTaskHistoryOperations(ctx, dbExecutor, taskID)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get task history operations: %w", err)
	}
	return operations, nil
}
