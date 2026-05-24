package projecttypehandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func addRequestToDomain(request addRequest) domain.ProjectType {
	return domain.ProjectType{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectType {
	return domain.ProjectType{
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func DomainToResponse(projectType domain.ProjectType) ProjectTypeResponse {
	return ProjectTypeResponse{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func domainArrayToResponseArray(projectTypes []domain.ProjectType) []ProjectTypeResponse {
	projectTypeResponses := []ProjectTypeResponse{}
	for _, projectType := range projectTypes {
		projectTypeResponses = append(projectTypeResponses, DomainToResponse(projectType))
	}
	return projectTypeResponses
}

func toSearchResponse(projectTypes []domain.ProjectType, pager browser.Result) searchResponse {
	return searchResponse{
		ProjectTypes: domainArrayToResponseArray(projectTypes),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
