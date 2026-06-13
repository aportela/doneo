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
	Add(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	Search(ctx context.Context, projectId string) ([]domain.HistoryOperation, error)
}

type historyOperationService struct {
	database   database.Database
	repository historyoperationrepository.HistoryOperationRepository
}

func NewService(database database.Database, repository historyoperationrepository.HistoryOperationRepository) HistoryOperationService {
	return &historyOperationService{database: database, repository: repository}
}

func (service *historyOperationService) Add(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddProjectOperation(ctx, projectId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *historyOperationService) Search(ctx context.Context, projectId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[HistoryOperationService] failed to get project history operations: %w", err)
	}
	return operations, nil
}
