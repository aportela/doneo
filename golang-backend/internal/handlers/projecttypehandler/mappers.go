package projecttypehandler

import (
	"github.com/aportela/doneo/internal/domain"
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

func domainToResponse(projectType domain.ProjectType) projectTypeResponse {
	return projectTypeResponse{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func domainArrayToResponseArray(projectTypes []domain.ProjectType) []projectTypeResponse {
	projectTypeResponses := []projectTypeResponse{}
	for _, projectType := range projectTypes {
		projectTypeResponses = append(projectTypeResponses, domainToResponse(projectType))
	}
	return projectTypeResponses
}

func toSearchResponse(users []domain.ProjectType) searchResponse {
	return searchResponse{
		ProjectTypes: domainArrayToResponseArray(users),
	}
}
