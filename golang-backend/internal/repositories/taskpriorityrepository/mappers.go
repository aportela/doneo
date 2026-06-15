package taskpriorityrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(taskPriority domain.TaskPriority) taskPriorityDTO {
	return taskPriorityDTO{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
		Index:    taskPriority.Index,
	}
}

func toDomain(taskPriority taskPriorityDTO) domain.TaskPriority {
	return domain.TaskPriority{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
		Index:    taskPriority.Index,
	}
}

func toDomainArray(taskPriorities []taskPriorityDTO) []domain.TaskPriority {
	results := make([]domain.TaskPriority, 0, len(taskPriorities))
	for _, projectPriority := range taskPriorities {
		results = append(results, toDomain(projectPriority))
	}
	return results
}

func toFilterDTO(filter domain.SearchTaskPrioritiesFilter) searchFilterDTO {
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
		if strings.Contains(sqlErr.Error(), "task_priorities.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		}
	}
	return err
}
