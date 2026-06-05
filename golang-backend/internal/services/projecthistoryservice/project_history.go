package projecthistoryservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	projecthistoryrepository "github.com/aportela/doneo/internal/repositories/projecthistoryrepository"
)

type ProjectHistoryService interface {
	Add(ctx context.Context, projectId string, operationType uint, operationDate int64, operationUserId string) error
	Search(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error)
}

type projectHistoryService struct {
	database   database.Database
	repository projecthistoryrepository.ProjectHistoryRepository
}

func NewService(database database.Database, repository projecthistoryrepository.ProjectHistoryRepository) ProjectHistoryService {
	return &projectHistoryService{database: database, repository: repository}
}

func (service *projectHistoryService) Add(ctx context.Context, projectId string, operationType uint, operationDate int64, operationUserId string) error {
	return service.repository.Add(ctx, projectId, operationType, operationDate, operationUserId)
}

func (service *projectHistoryService) Search(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error) {
	operations, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectHistoryService] failed to get project history operations: %w", err)
	}
	return operations, nil
}
