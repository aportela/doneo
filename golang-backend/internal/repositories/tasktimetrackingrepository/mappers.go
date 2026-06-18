package tasktimetrackingrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(taskTimeTracking domain.TaskTimeTracking) taskTimeTrackingDTO {
	return taskTimeTrackingDTO{
		ID:           taskTimeTracking.ID,
		CreatedAt:    taskTimeTracking.CreatedAt.UnixMilli(),
		CreatorID:    taskTimeTracking.CreatedBy.ID,
		CreatorName:  taskTimeTracking.CreatedBy.Name,
		Summary:      taskTimeTracking.Summary,
		TotalSeconds: taskTimeTracking.TotalSeconds,
	}
}

func toDomain(taskTimeTracking taskTimeTrackingDTO) domain.TaskTimeTracking {
	return domain.TaskTimeTracking{
		ID:           taskTimeTracking.ID,
		CreatedAt:    time.UnixMilli(taskTimeTracking.CreatedAt),
		CreatedBy:    domain.UserBase{ID: taskTimeTracking.CreatorID, Name: taskTimeTracking.CreatorName},
		Summary:      taskTimeTracking.Summary,
		TotalSeconds: taskTimeTracking.TotalSeconds,
	}
}

func toDomainArray(taskTimeTrackings []taskTimeTrackingDTO) []domain.TaskTimeTracking {
	results := make([]domain.TaskTimeTracking, 0, len(taskTimeTrackings))
	for _, timer := range taskTimeTrackings {
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
