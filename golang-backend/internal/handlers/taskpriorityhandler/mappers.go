package taskpriorityhandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func addRequestToDomain(request addRequest) domain.TaskPriority {
	return domain.TaskPriority{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskPriority {
	return domain.TaskPriority{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func domainToResponse(taskPriority domain.TaskPriority) taskPriorityResponse {
	return taskPriorityResponse{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
	}
}

func domainArrayToResponseArray(taskPriorities []domain.TaskPriority) []taskPriorityResponse {
	taskPriorityResponses := []taskPriorityResponse{}
	for _, taskPriority := range taskPriorities {
		taskPriorityResponses = append(taskPriorityResponses, domainToResponse(taskPriority))
	}
	return taskPriorityResponses
}

func toSearchResponse(taskPriorities []domain.TaskPriority, pager browser.Result) searchResponse {
	return searchResponse{
		TaskPriorities: domainArrayToResponseArray(taskPriorities),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
