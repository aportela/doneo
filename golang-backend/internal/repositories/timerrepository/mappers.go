package timerrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func toDTO(timer domain.Timer) timerDTO {
	return timerDTO{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToSQLNullInt64(timer.FinishedAt),
	}
}

func toDomain(timer timerDTO) domain.Timer {
	return domain.Timer{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  time.UnixMilli(timer.StartedAt),
		FinishedAt: utils.SQLNullInt64ToTimePtr(timer.FinishedAt),
	}
}

func toDomainArray(timers []timerDTO) []domain.Timer {
	results := make([]domain.Timer, 0, len(timers))
	for _, timer := range timers {
		results = append(results, toDomain(timer))
	}
	return results
}
