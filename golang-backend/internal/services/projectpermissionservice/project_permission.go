package projectpermissionservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrepository"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPermissionService interface {
	Add(ctx context.Context, projectId string, permission domain.ProjectPermission) (domain.ProjectPermission, error)
	Delete(ctx context.Context, projectId string, permissionId string) error
	Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error)
	Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error)
}

type projectPermissionService struct {
	database   database.Database
	repository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(database database.Database, repository projectpermissionrepository.ProjectPermissionRepository) ProjectPermissionService {
	return &projectPermissionService{database: database, repository: repository}
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
		return domain.ProjectPermission{}, fmt.Errorf("user ID not found in context")
	}
	permission.ID = utils.UUID()
	err = service.repository.Add(ctx, projectId, permission)
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionAdded})
	if err != nil {
		return domain.ProjectPermission{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.ProjectPermission{}, err
	}
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
		return fmt.Errorf("user ID not found in context")
	}
	err = service.repository.Delete(ctx, projectId, permissionId)
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *projectPermissionService) Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error) {
	projectPermission, err := service.repository.Get(ctx, permissionId)
	if err != nil {
		return domain.ProjectPermission{}, fmt.Errorf("[ProjectTypeService] failed to get project permission: %w", err)
	}
	return projectPermission, nil
}

func (service *projectPermissionService) Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error) {
	projectPermissions, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to get project permissions: %w", err)
	}
	return projectPermissions, nil
}
