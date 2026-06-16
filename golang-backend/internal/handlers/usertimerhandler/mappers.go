package timerhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func DomainToResponse(timer domain.UserTimer) timerResponse {
	return timerResponse{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToInt64Ptr(timer.FinishedAt),
	}
}

func domainArrayToResponseArray(timers []domain.UserTimer) []timerResponse {
	timerResponses := []timerResponse{}
	for _, timer := range timers {
		timerResponses = append(timerResponses, DomainToResponse(timer))
	}
	return timerResponses
}

func toSearchResponse(timers []domain.UserTimer) searchResponse {
	return searchResponse{
		Timers: domainArrayToResponseArray(timers),
	}
}
