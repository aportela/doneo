package projectpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
)

type ProjectPriorityService interface {
	Add(ctx context.Context, projectPriority domain.ProjectPriority) error
	Update(ctx context.Context, projectPriority domain.ProjectPriority) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectPriority, error)
	Search(ctx context.Context) ([]domain.ProjectPriority, error)
}

type projectPriorityService struct {
	repository projectpriorityrepository.ProjectPriorityRepository
}

func NewProjectPriorityService(repository projectpriorityrepository.ProjectPriorityRepository) ProjectPriorityService {
	return &projectPriorityService{repository: repository}
}

func (s *projectPriorityService) Add(ctx context.Context, projectPriority domain.ProjectPriority) error {
	return s.repository.Add(ctx, projectpriorityrepository.DomainToDTO(projectPriority))
}

func (s *projectPriorityService) Update(ctx context.Context, projectPriority domain.ProjectPriority) error {
	return s.repository.Update(ctx, projectpriorityrepository.DomainToDTO(projectPriority))
}

func (s *projectPriorityService) Delete(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *projectPriorityService) Get(ctx context.Context, id string) (domain.ProjectPriority, error) {
	projectPriority, err := s.repository.Get(ctx, id)
	if err != nil {
		return projectpriorityrepository.DTOToDomain(projectPriority), fmt.Errorf("[ProjectPriorityService] failed to get project priority with ID %s: %w", id, err)
	}
	return projectpriorityrepository.DTOToDomain(projectPriority), nil
}

func (s *projectPriorityService) Search(ctx context.Context) ([]domain.ProjectPriority, error) {
	projectPriorities, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ProjectPriorityService] failed to search project priorities: %w", err)
	}
	return projectpriorityrepository.DTOArrayToDomainArray(projectPriorities), nil
}
