package projectstatushandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddProjectStatusRequestToProjectStatusDomain(request addProjectStatusRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		Index:    request.Index,
		HexColor: request.HexColor,
	}
}

func mapUpdateProjectStatusRequestToProjectStatusDomain(request updateProjectStatusRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		Index:    request.Index,
		HexColor: request.HexColor,
	}
}

func mapProjectStatusDomainToProjectStatusResponse(projectStatus domain.ProjectStatus) ProjectStatusResponse {
	return ProjectStatusResponse{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		Index:    projectStatus.Index,
		HexColor: projectStatus.HexColor,
	}
}

func mapProjectStatusDomainToAddProjectStatusResponse(projectStatus domain.ProjectStatus) addProjectStatusResponse {
	return addProjectStatusResponse{
		ProjectStatus: mapProjectStatusDomainToProjectStatusResponse(projectStatus),
	}
}

func mapProjectStatusDomainToUpdateProjectStatusResponse(projectStatus domain.ProjectStatus) updateProjectStatusResponse {
	return updateProjectStatusResponse{
		ProjectStatus: mapProjectStatusDomainToProjectStatusResponse(projectStatus),
	}
}

func mapProjectStatusDomainToGetProjectStatusResponse(projectStatus domain.ProjectStatus) getProjectStatusResponse {
	return getProjectStatusResponse{
		ProjectStatus: mapProjectStatusDomainToProjectStatusResponse(projectStatus),
	}
}

func mapProjectStatusArrayDomainToProjectStatusArrayResponse(projectPriorities []domain.ProjectStatus) []ProjectStatusResponse {
	projectStatusResponses := []ProjectStatusResponse{}
	for _, projectStatus := range projectPriorities {
		projectStatusResponses = append(projectStatusResponses, mapProjectStatusDomainToProjectStatusResponse(projectStatus))
	}
	return projectStatusResponses
}

func mapProjectStatusArrayDomainToSearchProjectStatussResponse(users []domain.ProjectStatus) searchProjectStatusesResponse {
	return searchProjectStatusesResponse{
		ProjectStatuses: mapProjectStatusArrayDomainToProjectStatusArrayResponse(users),
	}
}

func mapSearchProjectPrioritysRequestToProjectPriorityFilterDomain(request searchProjectStatusesRequest) domain.ProjectStatusFilter {
	return domain.ProjectStatusFilter{
		WorkspaceId: request.WorkspaceId,
	}
}
