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

func DomainToResponse(projectPriority domain.ProjectPriority) ProjectPriorityResponse {
	return ProjectPriorityResponse{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func domainArrayToResponseArray(projectPriorities []domain.ProjectPriority) []ProjectPriorityResponse {
	projectPriorityResponses := []ProjectPriorityResponse{}
	for _, projectPriority := range projectPriorities {
		projectPriorityResponses = append(projectPriorityResponses, DomainToResponse(projectPriority))
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
