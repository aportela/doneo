package tasktimetrackingservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/tasktimetrackingrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type TaskTimeTrackingService interface {
	Add(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimeTracking) error
	Update(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimeTracking) error
	Delete(ctx context.Context, projectID string, taskID string, taskTimeEntryID string) error
	GetTaskTimerEntries(ctx context.Context, projectID string, taskID string) ([]domain.TaskTimeTracking, error)
}

type taskTimeTrackingService struct {
	db                         database.Database
	authorizationService       authorizationservice.AuthorizationService
	historyOperationService    historyoperationservice.HistoryOperationService
	taskTimeTrackingRepository tasktimetrackingrepository.TaskTimeTrackingRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, taskTimeTrackingRepository tasktimetrackingrepository.TaskTimeTrackingRepository) TaskTimeTrackingService {
	return &taskTimeTrackingService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, taskTimeTrackingRepository: taskTimeTrackingRepository}
}

func (service *taskTimeTrackingService) Add(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimeTracking) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		taskTimerEntry.CreatedBy.ID = contextUser.ID
		taskTimerEntry.CreatedBy.Name = contextUser.Name
		taskTimerEntry.CreatedAt = time.Now()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimeTrackingRepository.Add(ctx, tx, taskID, taskTimerEntry); err != nil {
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

func (service *taskTimeTrackingService) Update(ctx context.Context, projectID string, taskID string, taskTimerEntry domain.TaskTimeTracking) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimeTrackingRepository.Update(ctx, tx, taskTimerEntry); err != nil {
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

func (service *taskTimeTrackingService) Delete(ctx context.Context, projectID string, taskID string, taskTimeEntryID string) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.taskTimeTrackingRepository.Delete(ctx, tx, taskTimeEntryID); err != nil {
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

func (service *taskTimeTrackingService) GetTaskTimerEntries(ctx context.Context, projectID string, taskID string) ([]domain.TaskTimeTracking, error) {
	if _, err := service.authorizationService.RequireTaskViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if taskTimerEntries, err := service.taskTimeTrackingRepository.GetTaskTimerEntries(ctx, service.db, taskID); err != nil {
		return nil, fmt.Errorf("[TaskTimeTrackingService] failed to get task time tracking entries: %w", err)
	} else {
		return taskTimerEntries, nil
	}

}
