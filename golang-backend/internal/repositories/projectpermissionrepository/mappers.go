package projectpermissionrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(projectPermission domain.ProjectPermission) projectPermissionDTO {
	return projectPermissionDTO{
		ID:       projectPermission.ID,
		UserID:   projectPermission.User.ID,
		UserName: projectPermission.User.Name,
		RoleID:   projectPermission.Role.ID,
		RoleName: projectPermission.Role.Name,
	}
}

func toDomain(projectPermission projectPermissionDTO) domain.ProjectPermission {
	return domain.ProjectPermission{
		ID: projectPermission.ID,
		User: domain.UserBase{
			ID:   projectPermission.UserID,
			Name: projectPermission.UserName,
		},
		Role: domain.Role{
			RoleBase: domain.RoleBase{
				ID:   projectPermission.RoleID,
				Name: projectPermission.RoleName,
			},
			PermissionsBitmask: domain.Bitmask(projectPermission.RolePermissionsBitmask),
		},
	}
}

func toDomainArray(projectPermissions []projectPermissionDTO) []domain.ProjectPermission {
	results := make([]domain.ProjectPermission, 0, len(projectPermissions))
	for _, projectPermission := range projectPermissions {
		results = append(results, toDomain(projectPermission))
	}
	return results
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		if strings.Contains(sqlErr.Error(), "project_user_role.project_id, project_user_role.user_id") {
			return &domain.AlreadyExistsError{Field: "userId"}
		}
	}
	return err
}
