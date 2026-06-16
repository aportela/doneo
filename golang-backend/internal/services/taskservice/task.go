package taskservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/tagrepository"
	"github.com/aportela/doneo/internal/repositories/taskrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskService interface {
	Add(ctx context.Context, projectID string, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, projectID string, task domain.Task) (domain.Task, error)
	Delete(ctx context.Context, projectID string, taskID string) error
	Get(ctx context.Context, projectID string, taskID string) (domain.Task, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error)
}

type taskService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	taskRepository          taskrepository.TaskRepository
	tagRepository           tagrepository.TagRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, taskRepository taskrepository.TaskRepository, tagRepository tagrepository.TagRepository) TaskService {
	return &taskService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, taskRepository: taskRepository, tagRepository: tagRepository}
}

func (service *taskService) Add(ctx context.Context, projectID string, task domain.Task) (domain.Task, error) {
	err := service.authorizationService.WithTaskAddPermission(ctx, projectID, func(currentUserID string) error {
		task.ID = utils.UUID()
		task.CreatedBy.ID = currentUserID
		task.CreatedAt = time.Now()
		newTaskIndex, err := service.taskRepository.GetNextTaskIndex(ctx, service.db, projectID)
		if err != nil {
			return err
		}
		task.Index = newTaskIndex
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskRepository.Add(ctx, tx, projectID, task); err != nil {
				return err
			}
			if len(task.Tags) > 0 {
				for _, taskTag := range task.Tags {
					if err := service.tagRepository.AddTaskTag(ctx, tx, task.ID, taskTag); err != nil {
						return err
					}
				}
			}
			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				task.ID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     task.CreatedAt,
					OperationType: domain.EventTaskCreated,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (service *taskService) Update(ctx context.Context, projectID string, task domain.Task) (domain.Task, error) {
	err := service.authorizationService.WithTaskAddPermission(ctx, projectID, func(currentUserID string) error {
		task.UpdatedAt = utils.CurrentTimePtr()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskRepository.Update(ctx, tx, task); err != nil {
				return err
			}
			if err := service.tagRepository.DeleteTaskTags(ctx, tx, task.ID); err != nil {
				return err
			}
			if len(task.Tags) > 0 {
				for _, taskTag := range task.Tags {
					if err := service.tagRepository.AddTaskTag(ctx, tx, task.ID, taskTag); err != nil {
						return err
					}
				}
			}
			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				task.ID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     *task.UpdatedAt,
					OperationType: domain.EventTaskUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (service *taskService) Delete(ctx context.Context, projectID string, taskID string) error {
	err := service.authorizationService.WithTaskAddPermission(ctx, projectID, func(currentUserID string) error {
		deletedAt := time.Now()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskRepository.Delete(ctx, tx, taskID, deletedAt.UnixMilli()); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     deletedAt,
					OperationType: domain.EventTaskUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	return err
}

func (service *taskService) Get(ctx context.Context, projectID string, taskID string) (domain.Task, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Task{}, fmt.Errorf("[TaskService] user not found in context")
	}
	if err := service.authorizationService.RequireTaskViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return domain.Task{}, err
	}
	task, err := service.taskRepository.Get(ctx, service.db, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("[TaskService] failed to get task with ID %s: %w", taskID, err)
	}
	taskTags, err := service.tagRepository.GetTaskTags(ctx, service.db, task.ID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("[TaskService] failed to get task with ID %s: %w", taskID, err)
	}
	task.Tags = taskTags
	return task, nil
}

func (service *taskService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("[TaskService] user not found in context")
	}
	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		// filter by tasks visible by current user when admin flag is not set
		filter.ViewByUserId = &currentContextUserID
	}
	tasks, pagerResult, err := service.taskRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskService] failed to search tasks: %w", err)
	}
	return tasks, pagerResult, nil
}
