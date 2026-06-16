package taskpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskPriorityService interface {
	Add(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error)
	Update(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error)
	Delete(ctx context.Context, taskPriorityID string) error
	Get(ctx context.Context, taskPriorityID string) (domain.TaskPriority, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error)
}

type taskPriorityService struct {
	db                     database.Database
	authorizationService   authorizationservice.AuthorizationService
	taskPriorityRepository taskpriorityrepository.TaskPriorityRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, taskPriorityRepository taskpriorityrepository.TaskPriorityRepository) TaskPriorityService {
	return &taskPriorityService{db: db, authorizationService: authorizationService, taskPriorityRepository: taskPriorityRepository}
}

func (service *taskPriorityService) Add(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.TaskPriority{}, err
	}
	taskPriority.ID = utils.UUID()
	if err := service.taskPriorityRepository.Add(ctx, service.db, taskPriority); err != nil {
		return domain.TaskPriority{}, err
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Update(ctx context.Context, taskPriority domain.TaskPriority) (domain.TaskPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.TaskPriority{}, err
	}
	if err := service.taskPriorityRepository.Update(ctx, service.db, taskPriority); err != nil {
		return domain.TaskPriority{}, err
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Delete(ctx context.Context, taskPriorityID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.taskPriorityRepository.Delete(ctx, service.db, taskPriorityID); err != nil {
		return err
	}
	return nil
}

func (service *taskPriorityService) Get(ctx context.Context, taskPriorityID string) (domain.TaskPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.TaskPriority{}, err
	}
	taskPriority, err := service.taskPriorityRepository.Get(ctx, service.db, taskPriorityID)
	if err != nil {
		return domain.TaskPriority{}, fmt.Errorf("[TaskPriorityService] failed to get task priority with ID %s: %w", taskPriorityID, err)
	}
	return taskPriority, nil
}

func (service *taskPriorityService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.Result{}, err
	}
	taskPriorities, pagerResult, err := service.taskPriorityRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskPriorityService] failed to search task priorities: %w", err)
	}
	return taskPriorities, pagerResult, nil
}
