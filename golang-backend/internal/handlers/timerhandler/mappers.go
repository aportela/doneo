package timerhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func DomainToResponse(timer domain.Timer) TimerResponse {
	return TimerResponse{
		ID:         timer.ID,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToInt64Ptr(timer.FinishedAt),
	}
}

func domainArrayToResponseArray(timers []domain.Timer) []TimerResponse {
	timerResponses := []TimerResponse{}
	for _, timer := range timers {
		timerResponses = append(timerResponses, DomainToResponse(timer))
	}
	return timerResponses
}

func toSearchResponse(timers []domain.Timer) searchResponse {
	return searchResponse{
		Timers: domainArrayToResponseArray(timers),
	}
}
