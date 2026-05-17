package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func RoleToDTO(role domain.Role) RoleDTO {
	return RoleDTO{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: uint8(role.PermissionsBitmask),
	}
}

func DTOToRole(role RoleDTO) domain.Role {
	return domain.Role{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: domain.PermissionsBitmask(role.PermissionsBitmask),
	}
}

func ToRoleArray(roles []RoleDTO) []domain.Role {
	results := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		results = append(results, DTOToRole(role))
	}
	return results
}

func SearchRolesFilterToDTO(filter domain.SearchRolesFilter) SearchRolesFilterDTO {
	return SearchRolesFilterDTO{
		Name:                        filter.Name,
		RequiredPermissionsBitmask:  (*uint8)(filter.RequiredPermissionsBitmask),
		ForbiddenPermissionsBitmask: (*uint8)(filter.ForbiddenPermissionsBitmask),
	}
}
