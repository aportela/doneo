package taskstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/utils"
)

type TaskStatusService interface {
	Add(ctx context.Context, projectStatus domain.TaskStatus) (domain.TaskStatus, error)
	Update(ctx context.Context, projectStatus domain.TaskStatus) (domain.TaskStatus, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.TaskStatus, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error)
}

type taskStatusService struct {
	database   database.Database
	repository taskstatusrepository.ProjectStatusRepository
}

func NewService(db database.Database, repository taskstatusrepository.ProjectStatusRepository) TaskStatusService {
	return &taskStatusService{database: db, repository: repository}
}

func (service *taskStatusService) Add(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error) {
	taskStatus.ID = utils.UUID()
	err := service.repository.Add(ctx, taskStatus)
	if err != nil {
		return domain.TaskStatus{}, err
	}
	return taskStatus, nil
}

func (service *taskStatusService) Update(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error) {
	err := service.repository.Update(ctx, taskStatus)
	if err != nil {
		return domain.TaskStatus{}, err
	}
	return taskStatus, nil
}

func (service *taskStatusService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *taskStatusService) Get(ctx context.Context, id string) (domain.TaskStatus, error) {
	taskStatus, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.TaskStatus{}, fmt.Errorf("[TaskStatusService] failed to get task status with ID %s: %w", id, err)
	}
	return taskStatus, nil
}

func (service *taskStatusService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error) {
	taskStatuses, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskStatusService] failed to search task statuses: %w", err)
	}
	return taskStatuses, pagerResult, nil
}
