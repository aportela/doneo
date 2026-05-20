package projectpriorityhandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func addRequestToDomain(request addRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func domainToResponse(projectPriority domain.ProjectPriority) projectPriorityResponse {
	return projectPriorityResponse{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func domainArrayToResponseArray(projectPriorities []domain.ProjectPriority) []projectPriorityResponse {
	projectPriorityResponses := []projectPriorityResponse{}
	for _, projectPriority := range projectPriorities {
		projectPriorityResponses = append(projectPriorityResponses, domainToResponse(projectPriority))
	}
	return projectPriorityResponses
}

func toSearchResponse(projectPriorities []domain.ProjectPriority, pager browser.Result) searchResponse {
	return searchResponse{
		ProjectPriorities: domainArrayToResponseArray(projectPriorities),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
