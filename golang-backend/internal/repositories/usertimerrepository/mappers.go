package usertimerrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(timer domain.UserTimer) userTimerDTO {
	return userTimerDTO{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToSQLNullInt64(timer.FinishedAt),
	}
}

func toDomain(timer userTimerDTO) domain.UserTimer {
	return domain.UserTimer{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  time.UnixMilli(timer.StartedAt),
		FinishedAt: utils.SQLNullInt64ToTimePtr(timer.FinishedAt),
	}
}

func toDomainArray(timers []userTimerDTO) []domain.UserTimer {
	results := make([]domain.UserTimer, 0, len(timers))
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
