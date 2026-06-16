package authorizationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
)

type AuthorizationService interface {
	requireUserPermission(ctx context.Context, userID string, permission domain.Bitmask) error
	RequireUserAdminPermission(ctx context.Context, userID string) error
	WithUserAdminPermission(ctx context.Context, action func(currentUserID string) error) error

	requireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error
	RequireProjectUpdatePermission(ctx context.Context, userID string, projectID string) error
	WithProjectUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error
	RequireProjectDeletePermission(ctx context.Context, userID string, projectID string) error
	WithProjectDeletePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error
	RequireProjectViewPermission(ctx context.Context, userID string, projectID string) error

	requireTaskAddPermission(ctx context.Context, userID string, projectID string) error
	WithTaskAddPermission(ctx context.Context, projectID string, action func(currentUserID string) error) error
	RequireTaskUpdatePermission(ctx context.Context, userID string, projectID string) error
	WithTaskUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error
	RequireTaskDeletePermission(ctx context.Context, userID string, projectID string) error
	WithTaskDeletePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error
	RequireTaskViewPermission(ctx context.Context, userID string, projectID string) error
}

type authorizationService struct {
	db                          database.Database
	permissionCache             cache.PermissionCache
	userRepository              userrepository.UserRepository
	projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(db database.Database, cache cache.PermissionCache, userRepository userrepository.UserRepository, projectPermissionRepository projectpermissionrepository.ProjectPermissionRepository) AuthorizationService {
	return &authorizationService{
		db:                          db,
		permissionCache:             cache,
		userRepository:              userRepository,
		projectPermissionRepository: projectPermissionRepository,
	}
}

func (service *authorizationService) requireUserPermission(ctx context.Context, userID string, permission domain.Bitmask) error {
	userPermissionsBitmask, ok := service.permissionCache.GetUser(userID)
	if !ok {
		user, err := service.userRepository.Get(ctx, service.db, userID)
		if err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get user permissions: %w", err)
		}
		userPermissionsBitmask = user.PermissionsBitmask
		service.permissionCache.SetUser(userID, userPermissionsBitmask)
	}
	if userPermissionsBitmask.HasFlag(permission) {
		return nil
	}
	return domain.AuthorizationError
}

func (service *authorizationService) RequireUserAdminPermission(ctx context.Context, userID string) error {
	return service.requireUserPermission(ctx, userID, domain.UserPermissionAdmin)
}

func (service *authorizationService) WithUserAdminPermission(ctx context.Context, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) requireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error {
	userPermissionsBitmask, ok := service.permissionCache.GetUser(userID)
	if !ok {
		user, err := service.userRepository.Get(ctx, service.db, userID)
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
		projectPermissions, err := service.projectPermissionRepository.GetProjectPermissions(ctx, service.db, projectID)
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
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionUpdateProject)
}

func (service *authorizationService) WithProjectUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.RequireProjectUpdatePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) RequireProjectDeletePermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionDeleteProject)
}

func (service *authorizationService) WithProjectDeletePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.RequireProjectDeletePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) RequireProjectViewPermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionViewProject)
}

func (service *authorizationService) requireTaskAddPermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionAddTask)
}

func (service *authorizationService) WithTaskAddPermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.requireTaskAddPermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) RequireTaskUpdatePermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionUpdateTask)
}

func (service *authorizationService) WithTaskUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.RequireTaskUpdatePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) RequireTaskDeletePermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionDeleteTask)
}

func (service *authorizationService) WithTaskDeletePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AuthorizationService] user not found in context")
	}

	if err := service.RequireTaskDeletePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *authorizationService) RequireTaskViewPermission(ctx context.Context, userID string, projectID string) error {
	return service.requireProjectPermission(ctx, userID, projectID, domain.PermissionViewTask)
}
