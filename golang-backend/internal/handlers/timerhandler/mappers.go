package timerhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func DomainToResponse(timer domain.Timer) timerResponse {
	return timerResponse{
		ID:         timer.ID,
		Summary:    timer.Summary,
		StartedAt:  timer.StartedAt.UnixMilli(),
		FinishedAt: utils.TimePtrToInt64Ptr(timer.FinishedAt),
	}
}

func domainArrayToResponseArray(timers []domain.Timer) []timerResponse {
	timerResponses := []timerResponse{}
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
