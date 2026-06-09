package projectpermissionrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(projectPermission domain.ProjectPermission) projectPermissionDTO {
	return projectPermissionDTO{
		ID:       projectPermission.ID,
		UserId:   projectPermission.User.ID,
		UserName: projectPermission.User.Name,
		RoleId:   projectPermission.Role.ID,
		RoleName: projectPermission.Role.Name,
	}
}

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
