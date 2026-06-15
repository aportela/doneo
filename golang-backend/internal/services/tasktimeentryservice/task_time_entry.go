package tasktimeentryservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/tasktimeentryrepository"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskTimeEntryService interface {
	Add(ctx context.Context, projectId string, taskId string, taskTimeEntry domain.TaskTimeEntry) error
	Update(ctx context.Context, projectId string, taskId string, taskTimeEntry domain.TaskTimeEntry) error
	Delete(ctx context.Context, projectId string, taskId string, taskTimeEntryId string) error
	Get(ctx context.Context, id string) (domain.TaskTimeEntry, error)
	GetTaskTimeEntries(ctx context.Context, taskId string) ([]domain.TaskTimeEntry, error)
}

type taskTimeEntryService struct {
	database   database.Database
	history    historyoperationservice.HistoryOperationService
	repository tasktimeentryrepository.TaskTimeEntryRepository
}

func NewService(db database.Database, history historyoperationservice.HistoryOperationService, repository tasktimeentryrepository.TaskTimeEntryRepository) TaskTimeEntryService {
	return &taskTimeEntryService{database: db, history: history, repository: repository}
}

func (service *taskTimeEntryService) Add(ctx context.Context, projectId string, taskId string, taskTimeEntry domain.TaskTimeEntry) error {
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
		return fmt.Errorf("[TaskTimeEntryService] user ID not found in context")
	}
	taskTimeEntry.CreatedBy.ID = currentUserId
	taskTimeEntry.CreatedAt = time.Now()
	err = service.repository.Add(ctx, taskId, taskTimeEntry)
	if err != nil {
		return err
	}
	_, err = service.history.AddTaskHistoryOperation(ctx, tx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: taskTimeEntry.CreatedAt, OperationType: domain.EventTaskTimeEntryAdded})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (service *taskTimeEntryService) Update(ctx context.Context, projectId string, taskId string, taskTimeEntry domain.TaskTimeEntry) error {
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
		return fmt.Errorf("[TaskTimeEntryService] user ID not found in context")
	}
	if err := service.repository.Update(ctx, taskTimeEntry); err != nil {
		return fmt.Errorf("[TaskTimeEntryService] failed to update task time entry with ID %s: %w", taskTimeEntry.ID, err)
	}
	_, err = service.history.AddTaskHistoryOperation(ctx, tx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskTimeEntryUpdated})
	if err != nil {
		return err
	}
	return nil
}

func (service *taskTimeEntryService) Delete(ctx context.Context, projectId string, taskId string, taskTimeEntryId string) error {
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
		return fmt.Errorf("[TaskTimeEntryService] user ID not found in context")
	}
	if err := service.repository.Delete(ctx, taskTimeEntryId); err != nil {
		return fmt.Errorf("[TaskTimeEntryService] failed to delete task time entry with ID %s: %w", taskTimeEntryId, err)
	}
	_, err = service.history.AddTaskHistoryOperation(ctx, tx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskTimeEntryDeleted})
	if err != nil {
		return err
	}
	return nil
}

func (service *taskTimeEntryService) Get(ctx context.Context, id string) (domain.TaskTimeEntry, error) {
	taskTimeEntry, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.TaskTimeEntry{}, fmt.Errorf("[TaskTimeEntryService] failed to get task time entry with ID %s: %w", id, err)
	}
	return taskTimeEntry, nil
}

func (service *taskTimeEntryService) GetTaskTimeEntries(ctx context.Context, taskId string) ([]domain.TaskTimeEntry, error) {
	taskTimeEntries, err := service.repository.GetTaskTimeEntries(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("[TaskTimeEntryService] failed to get task time entries: %w", err)
	}
	return taskTimeEntries, nil
}
