package tasktimerentryservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
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
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		taskTimerEntry.CreatedBy.ID = contextUser.ID
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
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     taskTimerEntry.CreatedAt,
					OperationType: domain.EventTaskTimeEntryAdded,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *taskTimerEntryService) Update(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimerEntry) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimerEntryRepository.Update(ctx, tx, taskTimerEntry); err != nil {
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
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskTimeEntryUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *taskTimerEntryService) Delete(ctx context.Context, projectID string, taskID string, taskTimeEntryID string) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
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
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskTimeEntryDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *taskTimerEntryService) GetTaskTimerEntries(ctx context.Context, projectID string, taskID string) ([]domain.TaskTimerEntry, error) {
	if _, err := service.authorizationService.RequireTaskViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if taskTimerEntries, err := service.taskTimerEntryRepository.GetTaskTimerEntries(ctx, service.db, taskID); err != nil {
		return nil, fmt.Errorf("[TaskTimerEntryService] failed to get task timer entries: %w", err)
	} else {
		return taskTimerEntries, nil
	}

}
