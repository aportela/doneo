package projectstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
)

type ProjectStatusService interface {
	AddProjectStatus(ctx context.Context, projectStatus domain.ProjectStatus) error
	UpdateProjectStatus(ctx context.Context, projectStatus domain.ProjectStatus) error
	DeleteProjectStatus(ctx context.Context, id string) error
	GetProjectStatus(ctx context.Context, id string) (domain.ProjectStatus, error)
	SearchProjectStatuss(ctx context.Context) ([]domain.ProjectStatus, error)
}

type projectStatusService struct {
	repository projectstatusrepository.ProjectStatusRepository
}

func NewProjectStatusService(repository projectstatusrepository.ProjectStatusRepository) ProjectStatusService {
	return &projectStatusService{repository: repository}
}

func (s *projectStatusService) AddProjectStatus(ctx context.Context, projectStatus domain.ProjectStatus) error {
	return s.repository.Add(ctx, projectstatusrepository.MapProyectStatusDomainToProyectStatusDTO(projectStatus))
}

func (s *projectStatusService) UpdateProjectStatus(ctx context.Context, projectStatus domain.ProjectStatus) error {
	return s.repository.Update(ctx, projectstatusrepository.MapProyectStatusDomainToProyectStatusDTO(projectStatus))
}

func (s *projectStatusService) DeleteProjectStatus(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectStatusService) GetProjectStatus(ctx context.Context, id string) (domain.ProjectStatus, error) {
	projectStatus, err := s.repository.Get(ctx, id)
	if err != nil {
		return projectstatusrepository.MapProjectStatusDTOToProjectStatusDomain(projectStatus), fmt.Errorf("[ProjectStatusService] failed to get project status with ID %s: %w", id, err)
	}
	return projectstatusrepository.MapProjectStatusDTOToProjectStatusDomain(projectStatus), nil
}

func (s *projectStatusService) SearchProjectStatuss(ctx context.Context) ([]domain.ProjectStatus, error) {
	projectStatuss, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectStatusService] failed to search project statuses: %w", err)
	}
	return projectstatusrepository.MapProjectStatusArrayDTOToProjectStatusArrayDomain(projectStatuss), nil
}
