package usertimerhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func DomainToResponse(timer domain.UserTimer) userTimerResponse {
	return userTimerResponse{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToInt64Ptr(timer.FinishedAt),
	}
}

func domainArrayToResponseArray(timers []domain.UserTimer) []userTimerResponse {
	timerResponses := []userTimerResponse{}
	for _, timer := range timers {
		timerResponses = append(timerResponses, DomainToResponse(timer))
	}
	return timerResponses
}

func toSearchResponse(timers []domain.UserTimer) searchResponse {
	return searchResponse{
		UserTimers: domainArrayToResponseArray(timers),
	}
}
