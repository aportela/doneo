package projectpermissionservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
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
	return service.repository.Add(ctx, permissionId, projectId, userId, roleId)
}

func (service *projectPermissionService) Delete(ctx context.Context, projectId string, permissionId string) error {
	return service.repository.Delete(ctx, projectId, permissionId)
}

func (service *projectPermissionService) Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error) {
	projectPermissions, err := service.repository.Search(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to get project permissions: %w", err)
	}
	return projectPermissions, nil
}
