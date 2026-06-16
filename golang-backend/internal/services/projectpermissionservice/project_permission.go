package projectpermissionservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPermissionService interface {
	Add(ctx context.Context, projectID string, permission domain.ProjectPermission) (domain.ProjectPermission, error)
	Delete(ctx context.Context, projectID string, permissionID string) error
	GetProjectPermissions(ctx context.Context, projectID string) ([]domain.ProjectPermission, error)
}

type projectPermissionService struct {
	db                          database.Database
	authorizationService        authorizationservice.AuthorizationService
	historyOperationService     historyoperationservice.HistoryOperationService
	projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository) ProjectPermissionService {
	return &projectPermissionService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, projectPermissionRepository: projectPermissionRepository}
}

func (service *projectPermissionService) withProjectUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[ProjectPermissionService] user not found in context")
	}

	if err := service.authorizationService.RequireProjectUpdatePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *projectPermissionService) Add(ctx context.Context, projectID string, permission domain.ProjectPermission) (domain.ProjectPermission, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		permission.ID = utils.UUID()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.projectPermissionRepository.Add(ctx, tx, projectID, permission); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectPermissionAdded,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	return permission, nil
}

func (service *projectPermissionService) Delete(ctx context.Context, projectID string, permissionID string) error {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.projectPermissionRepository.Delete(ctx, tx, permissionID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectPermissionDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	return err
}

func (service *projectPermissionService) GetProjectPermissions(ctx context.Context, projectID string) ([]domain.ProjectPermission, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[ProjectPermissionService] user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	projectPermissions, err := service.projectPermissionRepository.GetProjectPermissions(ctx, service.db, projectID)
	if err != nil {
		return nil, fmt.Errorf("[ProjectPermissionService] failed to get project permissions: %w", err)
	}
	return projectPermissions, nil
}
