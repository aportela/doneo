package projectpriorityrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(projectPriority domain.ProjectPriority) projectPriorityDTO {
	return projectPriorityDTO{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
		Index:    projectPriority.Index,
	}
}

func toDomain(projectPriority projectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
		Index:    projectPriority.Index,
	}
}

func toDomainArray(projectPriorities []projectPriorityDTO) []domain.ProjectPriority {
	results := make([]domain.ProjectPriority, 0, len(projectPriorities))
	for _, projectPriority := range projectPriorities {
		results = append(results, toDomain(projectPriority))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectPrioritiesFilter) searchFilterDTO {
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
		if strings.Contains(sqlErr.Error(), "project_priorities.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
