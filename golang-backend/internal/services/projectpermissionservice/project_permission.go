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
)

type ProjectPermissionService interface {
	Add(ctx context.Context, permissionId string, projectId string, userId string, roleId string) error
	Delete(ctx context.Context, projectId string, permissionId string) error
	Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error)
}

type projectPermissionService struct {
	database   database.Database
	repository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(database database.Database, repository projectpermissionrepository.ProjectPermissionRepository) ProjectPermissionService {
	return &projectPermissionService{database: database, repository: repository}
}

func (service *projectPermissionService) Add(ctx context.Context, permissionId string, projectId string, userId string, roleId string) error {
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
	err = service.repository.Add(ctx, permissionId, projectId, userId, roleId)
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionAdded})
	return tx.Commit()
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
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectPermissionDeleted})
	return tx.Commit()
}

func (service *projectPermissionService) Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error) {
	projectPermissions, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to get project permissions: %w", err)
	}
	return projectPermissions, nil
}
