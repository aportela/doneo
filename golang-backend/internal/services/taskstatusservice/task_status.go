package taskstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskStatusService interface {
	Add(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error)
	Update(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error)
	Delete(ctx context.Context, taskStatusID string) error
	Get(ctx context.Context, taskStatusID string) (domain.TaskStatus, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error)
}

type taskStatusService struct {
	db                   database.Database
	authorizationService authorizationservice.AuthorizationService
	taskStatusRepository taskstatusrepository.TaskStatusRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, taskStatusRepository taskstatusrepository.TaskStatusRepository) TaskStatusService {
	return &taskStatusService{db: db, authorizationService: authorizationService, taskStatusRepository: taskStatusRepository}
}

func (service *taskStatusService) Add(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error) {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		taskStatus.ID = utils.UUID()
		err := service.taskStatusRepository.Add(ctx, service.db, taskStatus)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.TaskStatus{}, err
	}
	return taskStatus, nil
}

func (service *taskStatusService) Update(ctx context.Context, taskStatus domain.TaskStatus) (domain.TaskStatus, error) {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		err := service.taskStatusRepository.Update(ctx, service.db, taskStatus)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.TaskStatus{}, err
	}
	return taskStatus, nil
}

func (service *taskStatusService) Delete(ctx context.Context, taskStatusID string) error {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		err := service.taskStatusRepository.Delete(ctx, service.db, taskStatusID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (service *taskStatusService) Get(ctx context.Context, taskStatusID string) (domain.TaskStatus, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.TaskStatus{}, fmt.Errorf("[TaskStatusService] user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return domain.TaskStatus{}, err
	}
	taskStatus, err := service.taskStatusRepository.Get(ctx, service.db, taskStatusID)
	if err != nil {
		return domain.TaskStatus{}, fmt.Errorf("[TaskStatusService] failed to get task status with ID %s: %w", taskStatusID, err)
	}
	return taskStatus, nil
}

func (service *taskStatusService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("[TaskStatusService] user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return nil, browser.Result{}, err
	}
	taskStatuses, pagerResult, err := service.taskStatusRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskStatusService] failed to search task statuses: %w", err)
	}
	return taskStatuses, pagerResult, nil
}
