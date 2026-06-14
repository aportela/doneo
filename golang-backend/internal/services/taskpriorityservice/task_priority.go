package taskpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/utils"
)

type TaskPriorityService interface {
	Add(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error)
	Update(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.TaskPriority, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error)
}

type taskPriorityService struct {
	database   database.Database
	repository taskpriorityrepository.TaskPriorityRepository
}

func NewService(db database.Database, repository taskpriorityrepository.TaskPriorityRepository) TaskPriorityService {
	return &taskPriorityService{database: db, repository: repository}
}

func (service *taskPriorityService) Add(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error) {
	taskPriority.ID = utils.UUID()
	err := service.repository.Add(ctx, taskPriority)
	if err != nil {
		return domain.TaskPriority{}, err
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Update(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error) {
	err := service.repository.Update(ctx, taskPriority)
	if err != nil {
		return domain.TaskPriority{}, err
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *taskPriorityService) Get(ctx context.Context, id string) (domain.TaskPriority, error) {
	taskPriority, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.TaskPriority{}, fmt.Errorf("[TaskPriorityService] failed to get task priority with ID %s: %w", id, err)
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error) {
	taskPriorities, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskPriorityService] failed to search task priorities: %w", err)
	}
	return taskPriorities, pagerResult, nil
}
