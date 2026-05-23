package projectstatushandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func addRequestToDomain(request addRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func domainToResponse(projectStatus domain.ProjectStatus) projectStatusResponse {
	return projectStatusResponse{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		HexColor: projectStatus.HexColor,
	}
}

func domainArrayToResponseArray(projectStatuses []domain.ProjectStatus) []projectStatusResponse {
	projectStatusResponses := []projectStatusResponse{}
	for _, projectStatus := range projectStatuses {
		projectStatusResponses = append(projectStatusResponses, domainToResponse(projectStatus))
	}
	return projectStatusResponses
}

func toSearchResponse(projectStatuses []domain.ProjectStatus, pager browser.Result) searchResponse {
	return searchResponse{
		ProjectStatuses: domainArrayToResponseArray(projectStatuses),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
