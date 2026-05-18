package projectstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
)

type ProjectStatusService interface {
	Add(ctx context.Context, projectStatus domain.ProjectStatus) error
	Update(ctx context.Context, projectStatus domain.ProjectStatus) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectStatus, error)
	Search(ctx context.Context) ([]domain.ProjectStatus, error)
}

type projectStatusService struct {
	repository projectstatusrepository.ProjectStatusRepository
}

func NewProjectStatusService(repository projectstatusrepository.ProjectStatusRepository) ProjectStatusService {
	return &projectStatusService{repository: repository}
}

func (s *projectStatusService) Add(ctx context.Context, projectStatus domain.ProjectStatus) error {
	return s.repository.Add(ctx, projectstatusrepository.DomainToDTO(projectStatus))
}

func (s *projectStatusService) Update(ctx context.Context, projectStatus domain.ProjectStatus) error {
	return s.repository.Update(ctx, projectstatusrepository.DomainToDTO(projectStatus))
}

func (s *projectStatusService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectStatusService) Get(ctx context.Context, id string) (domain.ProjectStatus, error) {
	projectStatus, err := s.repository.Get(ctx, id)
	if err != nil {
		return projectstatusrepository.DTOToDomain(projectStatus), fmt.Errorf("[ProjectStatusService] failed to get project status with ID %s: %w", id, err)
	}
	return projectstatusrepository.DTOToDomain(projectStatus), nil
}

func (s *projectStatusService) Search(ctx context.Context) ([]domain.ProjectStatus, error) {
	projectStatuses, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectStatusService] failed to search project statuses: %w", err)
	}
	return projectstatusrepository.DTOArrayToDomainArray(projectStatuses), nil
}
