package tasktimerentryservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/tasktimerentryrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskTimerEntryService interface {
	Add(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimerEntry) error
	Update(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimerEntry) error
	Delete(ctx context.Context, projectID string, taskID string, taskTimeEntryID string) error
	GetTaskTimerEntries(ctx context.Context, projectID string, taskID string) ([]domain.TaskTimerEntry, error)
}

type taskTimerEntryService struct {
	db                       database.Database
	authorizationService     authorizationservice.AuthorizationService
	historyOperationService  historyoperationservice.HistoryOperationService
	taskTimerEntryRepository tasktimerentryrepository.TaskTimerEntryRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, taskTimerEntryRepository tasktimerentryrepository.TaskTimerEntryRepository) TaskTimerEntryService {
	return &taskTimerEntryService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, taskTimerEntryRepository: taskTimerEntryRepository}
}

func (service *taskTimerEntryService) Add(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimerEntry) error {
	err := service.authorizationService.WithTaskUpdatePermission(ctx, projectID, func(currentUserID string) error {
		taskTimerEntry.CreatedBy.ID = currentUserID
		taskTimerEntry.CreatedAt = time.Now()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimerEntryRepository.Add(ctx, tx, taskID, taskTimerEntry); err != nil {
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
					CreatedAt:     taskTimerEntry.CreatedAt,
					OperationType: domain.EventTaskTimeEntryAdded,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})
	return err
}

func (service *taskTimerEntryService) Update(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimerEntry) error {
	err := service.authorizationService.WithTaskUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimerEntryRepository.Add(ctx, tx, taskID, taskTimerEntry); err != nil {
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
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskTimeEntryUpdated,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})
	return err
}

func (service *taskTimerEntryService) Delete(ctx context.Context, projectID string, taskID string, taskTimeEntryID string) error {
	err := service.authorizationService.WithTaskUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimerEntryRepository.Delete(ctx, tx, taskTimeEntryID); err != nil {
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
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskTimeEntryDeleted,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})
	return err
}

func (service *taskTimerEntryService) GetTaskTimerEntries(ctx context.Context, projectID string, taskID string) ([]domain.TaskTimerEntry, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireTaskViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	taskTimerEntries, err := service.taskTimerEntryRepository.GetTaskTimerEntries(ctx, service.db, taskID)
	if err != nil {
		return nil, fmt.Errorf("[TaskTimeEntryService] failed to get task timer entries: %w", err)
	}
	return taskTimerEntries, nil
}
