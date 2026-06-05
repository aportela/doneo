package projectpermissionrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDomain(projectPermission projectPermissionDTO) domain.ProjectPermission {
	return domain.ProjectPermission{
		ID: projectPermission.ID,
		User: domain.UserBase{
			ID:   projectPermission.UserId,
			Name: projectPermission.UserName,
		},
		Role: domain.Role{
			RoleBase: domain.RoleBase{
				ID:   projectPermission.RoleId,
				Name: projectPermission.RoleName,
			},
			PermissionsBitmask: domain.PermissionsBitmask(projectPermission.PermissionsBitmask),
		},
	}
}

func toDomainArray(projectPermissions []projectPermissionDTO) []domain.ProjectPermission {
	results := make([]domain.ProjectPermission, 0, len(projectPermissions))
	for _, projectType := range projectPermissions {
		results = append(results, toDomain(projectType))
	}
	return results
}
