package projectpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPriorityService interface {
	Add(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error)
	Update(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectPriority, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error)
}

type projectPriorityService struct {
	database   database.Database
	repository projectpriorityrepository.ProjectPriorityRepository
}

func NewService(database database.Database, repository projectpriorityrepository.ProjectPriorityRepository) ProjectPriorityService {
	return &projectPriorityService{database: database, repository: repository}
}

func (service *projectPriorityService) Add(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error) {
	projectPriority.ID = utils.UUID()
	err := service.repository.Add(ctx, projectPriority)
	if err != nil {
		return domain.ProjectPriority{}, err
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Update(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error) {
	err := service.repository.Update(ctx, projectPriority)
	if err != nil {
		return domain.ProjectPriority{}, err
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *projectPriorityService) Get(ctx context.Context, id string) (domain.ProjectPriority, error) {
	projectPriority, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.ProjectPriority{}, fmt.Errorf("[ProjectPriorityService] failed to get project priority with ID %s: %w", id, err)
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error) {
	projectPriorities, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectPriorityService] failed to search project priorities: %w", err)
	}
	return projectPriorities, pagerResult, nil
}
