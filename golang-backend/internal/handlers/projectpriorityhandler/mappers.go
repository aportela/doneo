package projectpriorityhandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func addRequestToDomain(request addRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func domainToResponse(projectPriority domain.ProjectPriority) projectPriorityResponse {
	return projectPriorityResponse{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
		Index:    projectPriority.Index,
	}
}

func domainArrayToResponseArray(projectPriorities []domain.ProjectPriority) []projectPriorityResponse {
	projectPriorityResponses := []projectPriorityResponse{}
	for _, projectPriority := range projectPriorities {
		projectPriorityResponses = append(projectPriorityResponses, domainToResponse(projectPriority))
	}
	return projectPriorityResponses
}

func toSearchResponse(projectPriorities []domain.ProjectPriority) searchResponse {
	return searchResponse{
		ProjectPriorities: domainArrayToResponseArray(projectPriorities),
	}
}
