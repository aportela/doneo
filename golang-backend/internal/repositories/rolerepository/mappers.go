package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func RoleToDTO(role domain.Role) RoleDTO {
	return RoleDTO{
		ID:                role.ID,
		Name:              role.Name,
		PermissionBitmask: uint8(role.Permission),
	}
}

func DTOToRole(role RoleDTO) domain.Role {
	return domain.Role{
		ID:         role.ID,
		Name:       role.Name,
		Permission: domain.PermissionBitmask(role.PermissionBitmask),
	}
}

func ToRoleArray(roles []RoleDTO) []domain.Role {
	results := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		results = append(results, DTOToRole(role))
	}
	return results
}
