package projecthistoryservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrespository"
)

type ProjectHistoryService interface {
	GetProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error)
}

type projectHistoryService struct {
	repository projecthistoryrespository.ProjectHistoryRepository
}

func NewService(repository projecthistoryrespository.ProjectHistoryRepository) ProjectHistoryService {
	return &projectHistoryService{repository: repository}
}

func (service *projectHistoryService) GetProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error) {
	operations, err := service.repository.GetProjectHistoryOperations(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectHistoryService] failed to get project history operations: %w", err)
	}
	return projecthistoryrespository.DTOArrayToDomainArray(operations), nil
}
