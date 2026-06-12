package projecthistoryservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectHistoryService interface {
	Add(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error)
	Search(ctx context.Context, projectId string) ([]domain.HistoryOperation, error)
}

type projectHistoryService struct {
	database   database.Database
	repository projecthistoryrepository.ProjectHistoryRepository
}

func NewService(database database.Database, repository projecthistoryrepository.ProjectHistoryRepository) ProjectHistoryService {
	return &projectHistoryService{database: database, repository: repository}
}

func (service *projectHistoryService) Add(ctx context.Context, projectId string, operation domain.HistoryOperation) (domain.HistoryOperation, error) {
	operation.ID = utils.UUID()
	err := service.repository.AddProjectOperation(ctx, projectId, operation)
	if err != nil {
		return domain.HistoryOperation{}, err
	}
	return operation, nil
}

func (service *projectHistoryService) Search(ctx context.Context, projectId string) ([]domain.HistoryOperation, error) {
	operations, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectHistoryService] failed to get project history operations: %w", err)
	}
	return operations, nil
}
