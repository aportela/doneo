package projectstatusrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(projectStatus domain.ProjectStatus) projectStatusDTO {
	return projectStatusDTO{
		ID:           projectStatus.ID,
		Name:         projectStatus.Name,
		HexColor:     projectStatus.HexColor,
		Index:        projectStatus.Index,
		FlagsBitmask: uint64(projectStatus.Flags),
	}
}

func toDomain(projectStatus projectStatusDTO) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		HexColor: projectStatus.HexColor,
		Index:    projectStatus.Index,
		Flags:    domain.Bitmask(projectStatus.FlagsBitmask),
	}
}

func toDomainArray(projectStatuses []projectStatusDTO) []domain.ProjectStatus {
	results := make([]domain.ProjectStatus, 0, len(projectStatuses))
	for _, projectStatus := range projectStatuses {
		results = append(results, toDomain(projectStatus))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectStatusesFilter) searchFilterDTO {
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
		if strings.Contains(sqlErr.Error(), "project_statuses.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
