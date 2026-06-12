package projecttaskservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projecttaskrepository"
	"github.com/aportela/doneo/internal/repositories/tagrepository"
	"github.com/aportela/doneo/internal/repositories/taskhistoryrepository"
	"github.com/aportela/doneo/internal/utils"
)

type TaskService interface {
	Add(ctx context.Context, projectId string, task domain.Task) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) (domain.Task, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Task, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error)
}

type taskService struct {
	database   database.Database
	repository projecttaskrepository.TaskRepository
}

func NewService(database database.Database, repository projecttaskrepository.TaskRepository) TaskService {
	return &taskService{database: database, repository: repository}
}

func (service *taskService) Add(ctx context.Context, projectId string, task domain.Task) (domain.Task, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Task{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Task{}, fmt.Errorf("[TaskService] user ID not found in context")
	}
	task.ID = utils.UUID()
	task.CreatedBy.ID = currentUserId
	task.CreatedAt = time.Now()
	task.Index, err = service.repository.GetNextTaskIndex(ctx, projectId)
	if err != nil {
		return domain.Task{}, err
	}
	err = service.repository.Add(ctx, projectId, task)
	if err != nil {
		return domain.Task{}, err
	}
	if len(task.Tags) > 0 {
		tagRepository := tagrepository.NewRepository(service.database)
		for _, tag := range task.Tags {
			err = tagRepository.AddTaskTag(ctx, task.ID, tag)
			if err != nil {
				return domain.Task{}, err
			}
		}
	}
	err = taskhistoryrepository.NewRepository(service.database).Add(ctx, task.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: task.CreatedAt, OperationType: domain.EventTaskCreated})
	if err != nil {
		return domain.Task{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (service *taskService) Update(ctx context.Context, task domain.Task) (domain.Task, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Task{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Task{}, fmt.Errorf("[TaskService] user ID not found in context")
	}
	task.UpdatedAt = utils.CurrentTimePtr()
	err = service.repository.Update(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	tagRepository := tagrepository.NewRepository(service.database)
	err = tagRepository.DeleteTaskTags(ctx, task.ID)
	if err != nil {
		return domain.Task{}, err
	}
	if len(task.Tags) > 0 {
		for _, tag := range task.Tags {
			err = tagRepository.AddTaskTag(ctx, task.ID, tag)
			if err != nil {
				return domain.Task{}, err
			}
		}
	}
	err = taskhistoryrepository.NewRepository(service.database).Add(ctx, task.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskUpdated})
	if err != nil {
		return domain.Task{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

func (service *taskService) Delete(ctx context.Context, id string) error {
	tx, err := service.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TaskService] user ID not found in context")
	}
	err = service.repository.Delete(ctx, id, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	err = taskhistoryrepository.NewRepository(service.database).Add(ctx, id, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *taskService) Get(ctx context.Context, id string) (domain.Task, error) {
	task, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.Task{}, fmt.Errorf("[TaskService] failed to get task with ID %s: %w", id, err)
	}
	tags, err := tagrepository.NewRepository(service.database).GetTaskTags(ctx, task.ID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("[TaskService] failed to get task with ID %s: %w", id, err)
	}
	task.Tags = tags
	return task, nil
}

func (service *taskService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error) {
	tasks, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[TaskService] failed to search tasks: %w", err)
	}
	return tasks, pagerResult, nil
}
