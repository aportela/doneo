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
	if contextUser, err := service.authorizationService.RequireTaskAddPermission(ctx, projectID); err != nil {
		return domain.Task{}, err
	} else {
		task.ID = utils.UUID()
		task.CreatedBy.ID = contextUser.ID
		task.CreatedAt = time.Now()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if newTaskIndex, err := service.taskRepository.GetNextTaskIndex(ctx, service.db, projectID); err != nil {
				return err
			} else {
				task.Index = newTaskIndex
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
						CreatedBy:     domain.UserBase{ID: contextUser.ID},
						CreatedAt:     task.CreatedAt,
						OperationType: domain.EventTaskCreated,
					},
				); err != nil {
					return err
				}
				return nil
			}
		}); err != nil {
			return domain.Task{}, err
		}
		return task, nil
	}
}

func (service *taskService) Update(ctx context.Context, projectID string, task domain.Task) (domain.Task, error) {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return domain.Task{}, err
	} else {
		task.UpdatedAt = utils.CurrentTimePtr()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     *task.UpdatedAt,
					OperationType: domain.EventTaskUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Task{}, err
		}
		return task, nil
	}
}

func (service *taskService) Delete(ctx context.Context, projectID string, taskID string) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
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
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     deletedAt,
					OperationType: domain.EventTaskUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *taskService) Get(ctx context.Context, projectID string, taskID string) (domain.Task, error) {
	if _, err := service.authorizationService.RequireTaskViewPermission(ctx, projectID); err != nil {
		return domain.Task{}, err
	}
	if task, err := service.taskRepository.Get(ctx, service.db, taskID); err != nil {
		return domain.Task{}, fmt.Errorf("[TaskService] failed to get task with ID %s: %w", taskID, err)
	} else {
		if taskTags, err := service.tagRepository.GetTaskTags(ctx, service.db, task.ID); err != nil {
			return domain.Task{}, fmt.Errorf("[TaskService] failed to get task tags with ID %s: %w", taskID, err)
		} else {
			task.Tags = taskTags
			return task, nil
		}
	}
}

func (service *taskService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("[TaskService] user not found in context")
	}
	contextUser, err := service.authorizationService.RequireUserAdminPermission(ctx)
	if err != nil {
		// filter by tasks visible by current user when admin flag is not set
		filter.ViewByUserId = &contextUser.ID
	}
	tasks, pagerResult, err := service.taskRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskService] failed to search tasks: %w", err)
	}
	return tasks, pagerResult, nil
}
