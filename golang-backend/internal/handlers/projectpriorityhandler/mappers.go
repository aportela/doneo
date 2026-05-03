package projectpriorityhandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddProjectPriorityRequestToProjectPriorityDomain(request addProjectPriorityRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:          request.ID,
		WorkspaceId: request.WorkspaceId,
		Name:        request.Name,
		Index:       request.Index,
		HexColor:    request.HexColor,
	}
}

func mapUpdateProjectPriorityRequestToProjectPriorityDomain(request updateProjectPriorityRequest) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:          request.ID,
		WorkspaceId: request.WorkspaceId,
		Name:        request.Name,
		Index:       request.Index,
		HexColor:    request.HexColor,
	}
}

func mapProjectPriorityDomainToProjectPriorityResponse(projectPriority domain.ProjectPriority) projectPriorityResponse {
	return projectPriorityResponse{
		ID:          projectPriority.ID,
		WorkspaceId: projectPriority.WorkspaceId,
		Name:        projectPriority.Name,
		Index:       projectPriority.Index,
		HexColor:    projectPriority.HexColor,
	}
}

func mapProjectPriorityDomainToAddProjectPriorityResponse(projectPriority domain.ProjectPriority) addProjectPriorityResponse {
	return addProjectPriorityResponse{
		ProjectPriority: mapProjectPriorityDomainToProjectPriorityResponse(projectPriority),
	}
}

func mapProjectPriorityDomainToUpdateProjectPriorityResponse(projectPriority domain.ProjectPriority) updateProjectPriorityResponse {
	return updateProjectPriorityResponse{
		ProjectPriority: mapProjectPriorityDomainToProjectPriorityResponse(projectPriority),
	}
}

func mapProjectPriorityDomainToGetProjectPriorityResponse(projectPriority domain.ProjectPriority) getProjectPriorityResponse {
	return getProjectPriorityResponse{
		ProjectPriority: mapProjectPriorityDomainToProjectPriorityResponse(projectPriority),
	}
}

func mapProjectPriorityArrayDomainToProjectPriorityArrayResponse(projectPriorities []domain.ProjectPriority) []projectPriorityResponse {
	projectPriorityResponses := []projectPriorityResponse{}
	for _, projectPriority := range projectPriorities {
		projectPriorityResponses = append(projectPriorityResponses, mapProjectPriorityDomainToProjectPriorityResponse(projectPriority))
	}
	return projectPriorityResponses
}

func mapProjectPriorityArrayDomainToSearchProjectPrioritysResponse(users []domain.ProjectPriority) searchProjectPrioritysResponse {
	return searchProjectPrioritysResponse{
		ProjectPriorities: mapProjectPriorityArrayDomainToProjectPriorityArrayResponse(users),
	}
}

func mapSearchProjectPrioritysRequestToProjectPriorityFilterDomain(request searchProjectPrioritysRequest) domain.ProjectPriorityFilter {
	return domain.ProjectPriorityFilter{
		WorkspaceId: request.WorkspaceId,
	}
}
