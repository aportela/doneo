package projectstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectStatusService interface {
	Add(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error)
	Update(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectStatus, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error)
}

type projectStatusService struct {
	database   database.Database
	repository projectstatusrepository.ProjectStatusRepository
}

func NewService(database database.Database, repository projectstatusrepository.ProjectStatusRepository) ProjectStatusService {
	return &projectStatusService{database: database, repository: repository}
}

func (service *projectStatusService) Add(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error) {
	projectStatus.ID = utils.UUID()
	err := service.repository.Add(ctx, projectStatus)
	if err != nil {
		return domain.ProjectStatus{}, err
	}
	return projectStatus, nil
}

func (service *projectStatusService) Update(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error) {
	err := service.repository.Update(ctx, projectStatus)
	if err != nil {
		return domain.ProjectStatus{}, err
	}
	return projectStatus, nil
}

func (service *projectStatusService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *projectStatusService) Get(ctx context.Context, id string) (domain.ProjectStatus, error) {
	projectStatus, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.ProjectStatus{}, fmt.Errorf("[ProjectStatusService] failed to get project status with ID %s: %w", id, err)
	}
	return projectStatus, nil
}

func (service *projectStatusService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error) {
	projectStatuses, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectStatusService] failed to search project statuses: %w", err)
	}
	return projectStatuses, pagerResult, nil
}
