package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
)

type ProjectTypeService interface {
	AddProjectType(ctx context.Context, projectType domain.ProjectType) error
	UpdateProjectType(ctx context.Context, projectType domain.ProjectType) error
	DeleteProjectType(ctx context.Context, id string) error
	GetProjectType(ctx context.Context, id string) (domain.ProjectType, error)
	SearchProjectTypes(ctx context.Context) ([]domain.ProjectType, error)
}

type projectTypeService struct {
	repository projecttyperepository.ProjectTypeRepository
}

func NewProjectTypeService(repository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{repository: repository}
}

func (s *projectTypeService) AddProjectType(ctx context.Context, projectType domain.ProjectType) error {
	return s.repository.Add(ctx, projecttyperepository.MapProjectTypeDomainToProjectTypeDTO(projectType))
}

func (s *projectTypeService) UpdateProjectType(ctx context.Context, projectType domain.ProjectType) error {
	return s.repository.Update(ctx, projecttyperepository.MapProjectTypeDomainToProjectTypeDTO(projectType))
}

func (s *projectTypeService) DeleteProjectType(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectTypeService) GetProjectType(ctx context.Context, id string) (domain.ProjectType, error) {
	projectType, err := s.repository.Get(ctx, id)
	if err != nil {
		return projecttyperepository.MapProjectTypeDTOToProjectTypeDomain(projectType), fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", id, err)
	}
	return projecttyperepository.MapProjectTypeDTOToProjectTypeDomain(projectType), nil
}

func (s *projectTypeService) SearchProjectTypes(ctx context.Context) ([]domain.ProjectType, error) {
	projectTypes, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	}
	return projecttyperepository.MapProjectTypeArrayDTOToProjectTypeArrayDomain(projectTypes), nil
}
