package projectpermissionservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPermissionService interface {
	Add(ctx context.Context, projectId string, permission domain.ProjectPermission) (domain.ProjectPermission, error)
	Delete(ctx context.Context, projectId string, permissionId string) error
	Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error)
	Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error)
}

type projectPermissionService struct {
	database                database.Database
	cache                   cache.PermissionCache
	historyOperationService historyoperationservice.HistoryOperationService
	repository              projectpermissionrepository.ProjectPermissionRepository
}

func NewService(db database.Database, cache cache.PermissionCache, historyOperationService historyoperationservice.HistoryOperationService, repository projectpermissionrepository.ProjectPermissionRepository) ProjectPermissionService {
	return &projectPermissionService{database: db, cache: cache, historyOperationService: historyOperationService, repository: repository}
}

func (service *projectPermissionService) Add(ctx context.Context, projectId string, permission domain.ProjectPermission) (domain.ProjectPermission, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.ProjectPermission{}, err
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
		return domain.ProjectPermission{}, fmt.Errorf("[ProjectPermissionService] user ID not found in context")
	}
	permission.ID = utils.UUID()
	err = service.repository.Add(ctx, projectId, permission)
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	_, err = service.historyOperationService.AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionAdded})
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	service.cache.SetProject(currentUserId, projectId, permission.Role.PermissionsBitmask)
	return permission, nil
}

func (service *projectPermissionService) Delete(ctx context.Context, projectId string, permissionId string) error {
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
		return fmt.Errorf("[ProjectPermissionService] user ID not found in context")
	}
	err = service.repository.Delete(ctx, projectId, permissionId)
	if err != nil {
		return err
	}
	_, err = service.historyOperationService.AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionDeleted})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	service.cache.DeleteProject(currentUserId, projectId)
	return nil
}

// TODO: remove ???
func (service *projectPermissionService) Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error) {
	projectPermission, err := service.repository.Get(ctx, permissionId)
	if err != nil {
		return domain.ProjectPermission{}, fmt.Errorf("[ProjectPermissionService] failed to get project permission: %w", err)
	}
	return projectPermission, nil
}

func (service *projectPermissionService) Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error) {
	projectPermissions, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectPermissionService] failed to get project permissions: %w", err)
	}
	return projectPermissions, nil
}
