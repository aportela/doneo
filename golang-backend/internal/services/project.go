package services

import (
	"context"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
)

type ProjectService struct {
	repository *repositories.ProjectRepository
}

func NewProjectService(repository *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{
		repository: repository,
	}
}

func (s *ProjectService) AddProject(ctx context.Context, project domain.Project) error {
	return s.repository.Add(ctx, project)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project domain.Project) error {
	return s.repository.Update(ctx, project)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *ProjectService) GetProject(ctx context.Context, id string) (*domain.Project, error) {
	return s.repository.Get(ctx, id)
}

func (s *ProjectService) SearchProjects(ctx context.Context) ([]domain.Project, error) {
	return s.repository.Search(ctx)
}
