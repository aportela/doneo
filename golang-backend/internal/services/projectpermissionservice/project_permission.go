package projectpermissionservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPermissionService interface {
	Add(ctx context.Context, projectID string, projectPermission domain.ProjectPermission) (domain.ProjectPermission, error)
	Delete(ctx context.Context, projectID string, projectPermissionID string) error
	GetProjectPermissions(ctx context.Context, projectID string) ([]domain.ProjectPermission, error)
}

type projectPermissionService struct {
	db                          database.Database
	permissionCache             cache.PermissionCache
	authorizationService        authorizationservice.AuthorizationService
	historyOperationService     historyoperationservice.HistoryOperationService
	projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(db database.Database, cache cache.PermissionCache, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository) ProjectPermissionService {
	return &projectPermissionService{db: db, permissionCache: cache, authorizationService: authorizationService, historyOperationService: historyOperationService, projectPermissionRepository: projectPermissionRepository}
}

func (service *projectPermissionService) Add(ctx context.Context, projectID string, projectPermission domain.ProjectPermission) (domain.ProjectPermission, error) {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return domain.ProjectPermission{}, err
	} else {
		projectPermission.ID = utils.UUID()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.projectPermissionRepository.Add(ctx, tx, projectID, projectPermission); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectPermissionAdded,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.ProjectPermission{}, err
		}
		if projectPermission, err := service.projectPermissionRepository.Get(ctx, service.db, projectPermission.ID); err != nil {
			return domain.ProjectPermission{}, err
		} else {
			return projectPermission, nil
		}
	}
}

func (service *projectPermissionService) Delete(ctx context.Context, projectID string, projectPermissionID string) error {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.projectPermissionRepository.Delete(ctx, tx, projectPermissionID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectPermissionDeleted,
				},
			); err != nil {
				return err
			}
			service.permissionCache.DeleteProject(contextUser.ID, projectID)
			return nil
		})
	}
}

func (service *projectPermissionService) GetProjectPermissions(ctx context.Context, projectID string) ([]domain.ProjectPermission, error) {
	if _, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return nil, err
	} else {
		if projectPermissions, err := service.projectPermissionRepository.GetProjectPermissions(ctx, service.db, projectID); err != nil {
			return nil, fmt.Errorf("[ProjectPermissionService] failed to get project permissions: %w", err)
		} else {
			return projectPermissions, nil
		}
	}
}
