package authorizationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
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
	database                    database.Database
	permissionCache             cache.PermissionCache
	userRepository              userrepository.UserRepository
	projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(database database.Database, cache cache.PermissionCache) AuthorizationService {
	return &authorizationService{
		database:                    database,
		permissionCache:             cache,
		userRepository:              userrepository.NewRepository(database),
		projectPermissionRepository: projectpermissionrepository.NewRepository(database),
	}
}

func (service *authorizationService) RequireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error {
	userPermissionsBitmask, ok := service.permissionCache.GetUser(userID)
	if !ok {
		user, err := service.userRepository.Get(ctx, userID)
		if err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get user permissions: %w", err)
		}
		userPermissionsBitmask = user.PermissionsBitmask
		service.permissionCache.SetUser(userID, userPermissionsBitmask)
	}
	if userPermissionsBitmask.HasFlag(domain.UserPermissionAdmin) {
		return nil
	}
	permissionsBitmask, ok := service.permissionCache.GetProject(userID, projectID)
	if !ok {
		found := false
		projectPermissions, err := service.projectPermissionRepository.Search(ctx, projectID)
		if err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get project permissions: %w", err)
		}
		for _, projectPermission := range projectPermissions {
			if projectPermission.User.ID == userID {
				permissionsBitmask = projectPermission.Role.PermissionsBitmask
				service.permissionCache.SetProject(userID, projectID, permissionsBitmask)
				found = true
			}
		}
		if !found {
			return domain.AuthorizationError
		}
	}
	if permissionsBitmask.HasFlag(permission) {
		return nil
	}
	return domain.AuthorizationError
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
