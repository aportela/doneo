package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func RoleToDTO(role domain.Role) RoleDTO {
	return RoleDTO{
		ID:                role.ID,
		Name:              role.Name,
		PermissionBitmask: uint8(role.PermissionsBitmask),
	}
}

func DTOToRole(role RoleDTO) domain.Role {
	return domain.Role{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: domain.PermissionsBitmask(role.PermissionBitmask),
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
		Name:              filter.Name,
		PermissionBitmask: (*uint8)(filter.PermissionsBitmask),
	}
}
