package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func RoleToDTO(role domain.Role) RoleDTO {
	return RoleDTO{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: uint64(role.PermissionsBitmask),
	}
}

func DTOToDomain(role RoleDTO) domain.Role {
	return domain.Role{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: domain.PermissionsBitmask(role.PermissionsBitmask),
	}
}

func DTOArrayToDomainArray(roles []RoleDTO) []domain.Role {
	results := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		results = append(results, DTOToDomain(role))
	}
	return results
}

func SearchRolesFilterToDTO(filter domain.SearchRolesFilter) SearchRolesFilterDTO {
	return SearchRolesFilterDTO{
		Name: filter.Name,
	}
}
