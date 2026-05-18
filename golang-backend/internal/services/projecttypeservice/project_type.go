package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
)

type ProjectTypeService interface {
	Add(ctx context.Context, projectType domain.ProjectType) error
	Update(ctx context.Context, projectType domain.ProjectType) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectType, error)
	Search(ctx context.Context) ([]domain.ProjectType, error)
}

type projectTypeService struct {
	repository projecttyperepository.ProjectTypeRepository
}

func NewProjectTypeService(repository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{repository: repository}
}

func (s *projectTypeService) Add(ctx context.Context, projectType domain.ProjectType) error {
	return s.repository.Add(ctx, projecttyperepository.DomainToDTO(projectType))
}

func (s *projectTypeService) Update(ctx context.Context, projectType domain.ProjectType) error {
	return s.repository.Update(ctx, projecttyperepository.DomainToDTO(projectType))
}

func (s *projectTypeService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectTypeService) Get(ctx context.Context, id string) (domain.ProjectType, error) {
	projectType, err := s.repository.Get(ctx, id)
	if err != nil {
		return projecttyperepository.DTOToDomain(projectType), fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", id, err)
	}
	return projecttyperepository.DTOToDomain(projectType), nil
}

func (s *projectTypeService) Search(ctx context.Context) ([]domain.ProjectType, error) {
	projectTypes, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	}
	return projecttyperepository.DTOArrayToDomainArray(projectTypes), nil
}
