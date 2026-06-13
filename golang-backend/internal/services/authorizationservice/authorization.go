package authorizationservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
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
	database   database.Database
	repository projectpermissionrepository.ProjectPermissionRepository
}

func NewService(database database.Database, repository projectpermissionrepository.ProjectPermissionRepository) AuthorizationService {
	return &authorizationService{database: database, repository: repository}
}

func (service *authorizationService) RequireProjectPermission(ctx context.Context, userID string, projectID string, permission domain.Bitmask) error {
	projectPermissions, err := service.repository.Search(ctx, projectID)
	if err != nil {
		return fmt.Errorf("[AuthorizationService] failed to get project permissions: %w", err)
	}
	for _, projectPermission := range projectPermissions {
		if projectPermission.User.ID == userID {
			if projectPermission.Role.PermissionsBitmask.HasFlag(permission) {
				return nil
			}
		}
	}
	return domain.AuthorizationError
	/*
		return fmt.Errorf(
			"[AuthorizationService] user %s lacks permission %v on project %s",
			userID, permission, projectID,
		)
	*/
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
