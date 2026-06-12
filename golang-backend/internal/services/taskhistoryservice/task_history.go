package taskhistoryservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskhistoryrepository"
	"github.com/aportela/doneo/internal/utils"
)

type TaskHistoryService interface {
	Add(ctx context.Context, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	Search(ctx context.Context, taskId string) ([]domain.HistoryOperation, error)
}

type taskHistoryService struct {
	database   database.Database
	repository taskhistoryrepository.TaskHistoryRepository
}

func NewService(database database.Database, repository taskhistoryrepository.TaskHistoryRepository) TaskHistoryService {
	return &taskHistoryService{database: database, repository: repository}
}

func (service *taskHistoryService) Add(ctx context.Context, taskId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.Add(ctx, taskId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *taskHistoryService) Search(ctx context.Context, taskId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.Search(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("[TaskHistoryService] failed to get task history operations: %w", err)
	}
	return operations, nil
}
