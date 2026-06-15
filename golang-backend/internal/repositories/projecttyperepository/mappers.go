package projecttyperepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(projectType domain.ProjectType) projectTypeDTO {
	return projectTypeDTO{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func toDomain(projectType projectTypeDTO) domain.ProjectType {
	return domain.ProjectType{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func toDomainArray(projectTypes []projectTypeDTO) []domain.ProjectType {
	results := make([]domain.ProjectType, 0, len(projectTypes))
	for _, projectType := range projectTypes {
		results = append(results, toDomain(projectType))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectTypesFilter) searchFilterDTO {
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
		if strings.Contains(sqlErr.Error(), "project_types.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
