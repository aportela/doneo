package taskstatusrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(taskStatus domain.TaskStatus) taskStatusDTO {
	return taskStatusDTO{
		ID:           taskStatus.ID,
		Name:         taskStatus.Name,
		HexColor:     taskStatus.HexColor,
		Index:        taskStatus.Index,
		FlagsBitmask: uint64(taskStatus.Flags),
	}
}

func toDomain(taskStatus taskStatusDTO) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       taskStatus.ID,
		Name:     taskStatus.Name,
		HexColor: taskStatus.HexColor,
		Index:    taskStatus.Index,
		Flags:    domain.Bitmask(taskStatus.FlagsBitmask),
	}
}

func toDomainArray(taskStatuses []taskStatusDTO) []domain.TaskStatus {
	results := make([]domain.TaskStatus, 0, len(taskStatuses))
	for _, projectStatus := range taskStatuses {
		results = append(results, toDomain(projectStatus))
	}
	return results
}

func toFilterDTO(filter domain.SearchTaskStatusesFilter) searchFilterDTO {
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
		if strings.Contains(sqlErr.Error(), "task_statuses.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
