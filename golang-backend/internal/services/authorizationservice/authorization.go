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
	requireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error

	RequireUserAdminPermission(ctx context.Context) (middlewares.ContextUser, error)

	RequireProjectUpdatePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
	RequireProjectDeletePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
	RequireProjectViewPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)

	RequireTaskAddPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
	RequireTaskUpdatePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
	RequireTaskDeletePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
	RequireTaskViewPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error)
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
		if user, err := service.userRepository.Get(ctx, service.db, userID); err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get user permissions: %w", err)
		} else {
			userPermissionsBitmask = user.PermissionsBitmask
			service.permissionCache.SetUser(userID, userPermissionsBitmask)
		}
	}
	if userPermissionsBitmask.HasFlag(permission) {
		return nil
	}
	return domain.AuthorizationError
}

func (service *authorizationService) requireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error {
	userPermissionsBitmask, ok := service.permissionCache.GetUser(userID)
	if !ok {
		if user, err := service.userRepository.Get(ctx, service.db, userID); err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get user permissions: %w", err)
		} else {
			userPermissionsBitmask = user.PermissionsBitmask
			service.permissionCache.SetUser(userID, userPermissionsBitmask)
		}
	}
	if userPermissionsBitmask.HasFlag(domain.UserPermissionAdmin) {
		return nil
	}
	projectPermissionsBitmask, ok := service.permissionCache.GetProject(userID, projectID)
	if !ok {
		projectPermissions, err := service.projectPermissionRepository.GetProjectPermissions(ctx, service.db, projectID)
		if err != nil {
			return fmt.Errorf("[AuthorizationService] failed to get project permissions: %w", err)
		}
		found := false
		for _, projectPermission := range projectPermissions {
			if projectPermission.User.ID == userID {
				projectPermissionsBitmask = projectPermission.Role.PermissionsBitmask
				service.permissionCache.SetProject(userID, projectID, projectPermissionsBitmask)
				found = true
			}
		}
		if !found {
			return domain.AuthorizationError
		}
	}
	if projectPermissionsBitmask.HasFlag(permission) {
		return nil
	}
	return domain.AuthorizationError
}

func (service *authorizationService) RequireUserAdminPermission(ctx context.Context) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireUserPermission(ctx, contextUser.ID, domain.UserPermissionAdmin); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireProjectUpdatePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionUpdateProject); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireProjectDeletePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionDeleteProject); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireProjectViewPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionViewProject); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireTaskAddPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionAddTask); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireTaskUpdatePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionUpdateTask); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireTaskDeletePermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionDeleteTask); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}

func (service *authorizationService) RequireTaskViewPermission(ctx context.Context, projectID string) (middlewares.ContextUser, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return middlewares.ContextUser{}, fmt.Errorf("[AuthorizationService] user not found in context")
	}
	if !contextUser.SkipAuthorization {
		if err := service.requireProjectPermission(ctx, contextUser.ID, projectID, domain.PermissionViewTask); err != nil {
			return middlewares.ContextUser{}, err
		}
	}
	return contextUser, nil
}
