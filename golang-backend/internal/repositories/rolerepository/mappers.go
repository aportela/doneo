package rolerepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		if strings.Contains(sqlErr.Error(), "roles.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
