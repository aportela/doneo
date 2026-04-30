package projectservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
)

type ProjectService interface {
	AddProject(ctx context.Context, Project domain.Project) error
	UpdateProject(ctx context.Context, Project domain.Project) error
	DeleteProject(ctx context.Context, id string) error
	GetProject(ctx context.Context, id string) (domain.Project, error)
	SearchProjects(ctx context.Context) ([]domain.Project, error)
}

type projectService struct {
	repository projectrepository.ProjectRepository
}

func NewProjectService(repository projectrepository.ProjectRepository) ProjectService {
	return &projectService{repository: repository}
}

func (s *projectService) AddProject(ctx context.Context, project domain.Project) error {
	return s.repository.Add(ctx, projectrepository.MapProyectDomainToProyectDTO(project))
}

func (s *projectService) UpdateProject(ctx context.Context, project domain.Project) error {
	return s.repository.Update(ctx, projectrepository.MapProyectDomainToProyectDTO(project))
}

func (s *projectService) DeleteProject(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectService) GetProject(ctx context.Context, id string) (domain.Project, error) {
	project, err := s.repository.Get(ctx, id)
	if err != nil {
		return projectrepository.MapProjectDTOToProjectDomain(project), fmt.Errorf("[ProjectService] failed to get project with ID %s: %w", id, err)
	}
	return projectrepository.MapProjectDTOToProjectDomain(project), nil
}

func (s *projectService) SearchProjects(ctx context.Context) ([]domain.Project, error) {
	projects, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectService] failed to search projects: %w", err)
	}
	return projectrepository.MapProjectArrayDTOToProjectArrayDomain(projects), nil
}
