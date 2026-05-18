package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func DomainToDTO(role domain.Role) roleDTO {
	return roleDTO{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: uint64(role.PermissionsBitmask),
	}
}

func DTOToDomain(role roleDTO) domain.Role {
	return domain.Role{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: domain.PermissionsBitmask(role.PermissionsBitmask),
	}
}

func DTOArrayToDomainArray(roles []roleDTO) []domain.Role {
	results := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		results = append(results, DTOToDomain(role))
	}
	return results
}

func DomainFilterToDTO(filter domain.SearchRolesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
