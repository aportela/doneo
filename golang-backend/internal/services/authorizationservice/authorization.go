package authorizationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/services/cacheservice"
)

type AuthorizationService interface {
	RequireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error
	RequireProjectUpdatePermission(ctx context.Context, userID string, projectID string) error
	RequireProjectDeletePermission(ctx context.Context, userID string, projectID string) error
	RequireProjectViewPermission(ctx context.Context, userID string, projectID string) error
	RequireTaskAddPermission(ctx context.Context, userID string, projectID string) error
	RequireTaskUpdatePermission(ctx context.Context, userID string, projectID string) error
	RequireTaskDeletePermission(ctx context.Context, userID string, projectID string) error
	RequireTaskViewPermission(ctx context.Context, userID string, projectID string) error
}

type authorizationService struct {
	database   database.Database // TODO: delete ||create project permission repository with this on NewService
	cache      cacheservice.PermissionCache
	repository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(database database.Database, repository projectpermissionrepository.ProjectPermissionRepository) AuthorizationService {
	return &authorizationService{database: database, cache: cacheservice.NewPermissionCache(), repository: repository}
}

func (service *authorizationService) RequireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error {
	permissionsBitmask, ok := service.cache.Get(userID, projectID)
	if !ok {
		found := false
		projectPermissions, err := service.repository.Search(ctx, projectID)
		if err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get project permissions: %w", err)
		}
		for _, projectPermission := range projectPermissions {
			if projectPermission.User.ID == userID {
				permissionsBitmask = projectPermission.Role.PermissionsBitmask
				service.cache.Set(userID, projectID, permissionsBitmask)
				found = true
			}
		}
		if !found {
			return domain.AuthorizationError
		}
	}
	if permissionsBitmask.HasFlag(permission) {
		return nil
	} else {
		return domain.AuthorizationError
	}
}

func (service *authorizationService) RequireProjectUpdatePermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionUpdateProject)
}

func (service *authorizationService) RequireProjectDeletePermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionDeleteProject)
}

func (service *authorizationService) RequireProjectViewPermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionViewProject)
}

func (service *authorizationService) RequireTaskAddPermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionAddTask)
}

func (service *authorizationService) RequireTaskUpdatePermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionUpdateTask)
}

func (service *authorizationService) RequireTaskDeletePermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionDeleteTask)
}

func (service *authorizationService) RequireTaskViewPermission(ctx context.Context, userID string, projectID string) error {
	return service.RequireProjectPermission(ctx, userID, projectID, domain.PermissionViewTask)
}
