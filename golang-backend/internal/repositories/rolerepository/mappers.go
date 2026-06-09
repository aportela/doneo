package rolerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(role domain.Role) roleDTO {
	return roleDTO{
		ID:                 role.ID,
		Name:               role.Name,
		PermissionsBitmask: uint64(role.PermissionsBitmask),
	}
}

func toDomain(role roleDTO) domain.Role {
	return domain.Role{
		RoleBase: domain.RoleBase{
			ID:   role.ID,
			Name: role.Name,
		},
		PermissionsBitmask: domain.Bitmask(role.PermissionsBitmask),
	}
}

func toDomainArray(roles []roleDTO) []domain.Role {
	results := make([]domain.Role, 0, len(roles))
	for _, role := range roles {
		results = append(results, toDomain(role))
	}
	return results
}

func toFilterDTO(filter domain.SearchRolesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
