package tasktimeentryrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(taskTimeEntry domain.TaskTimerEntry) taskTimerEntryDTO {
	return taskTimerEntryDTO{
		ID:           taskTimeEntry.ID,
		CreatedAt:    taskTimeEntry.CreatedAt.UnixMilli(),
		CreatorId:    taskTimeEntry.CreatedBy.ID,
		CreatorName:  taskTimeEntry.CreatedBy.Name,
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func toDomain(taskTimeEntry taskTimerEntryDTO) domain.TaskTimerEntry {
	return domain.TaskTimerEntry{
		ID:           taskTimeEntry.ID,
		CreatedAt:    time.UnixMilli(taskTimeEntry.CreatedAt),
		CreatedBy:    domain.UserBase{ID: taskTimeEntry.CreatorId, Name: taskTimeEntry.CreatorName},
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func toDomainArray(timers []taskTimerEntryDTO) []domain.TaskTimerEntry {
	results := make([]domain.TaskTimerEntry, 0, len(timers))
	for _, timer := range timers {
		results = append(results, toDomain(timer))
	}
	return results
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(summary)") {
			return &domain.ValidationError{Field: "summary"}
		}
	}
	return err
}
