package projectstatushandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func addRequestToDomain(request addRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func domainToResponse(projectStatus domain.ProjectStatus) projectStatusResponse {
	return projectStatusResponse{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		HexColor: projectStatus.HexColor,
		Index:    projectStatus.Index,
	}
}

func domainArrayToResponseArray(projectStatuses []domain.ProjectStatus) []projectStatusResponse {
	projectStatusResponses := []projectStatusResponse{}
	for _, projectStatus := range projectStatuses {
		projectStatusResponses = append(projectStatusResponses, domainToResponse(projectStatus))
	}
	return projectStatusResponses
}

func toSearchResponse(users []domain.ProjectStatus) searchResponse {
	return searchResponse{
		ProjectStatuses: domainArrayToResponseArray(users),
	}
}
