package projecttypehandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddProjectTypeRequestToProjectTypeDomain(request addProjectTypeRequest) domain.ProjectType {
	return domain.ProjectType{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func mapUpdateProjectTypeRequestToProjectTypeDomain(request updateProjectTypeRequest) domain.ProjectType {
	return domain.ProjectType{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
	}
}

func mapProjectTypeDomainToProjectTypeResponse(projectType domain.ProjectType) projectTypeResponse {
	return projectTypeResponse{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func mapProjectTypeDomainToAddProjectTypeResponse(projectType domain.ProjectType) addProjectTypeResponse {
	return addProjectTypeResponse{
		ProjectType: mapProjectTypeDomainToProjectTypeResponse(projectType),
	}
}

func mapProjectTypeDomainToUpdateProjectTypeResponse(projectType domain.ProjectType) updateProjectTypeResponse {
	return updateProjectTypeResponse{
		ProjectType: mapProjectTypeDomainToProjectTypeResponse(projectType),
	}
}

func mapProjectTypeDomainToGetProjectTypeResponse(projectType domain.ProjectType) getProjectTypeResponse {
	return getProjectTypeResponse{
		ProjectType: mapProjectTypeDomainToProjectTypeResponse(projectType),
	}
}

func mapProjectTypeArrayDomainToProjectTypeArrayResponse(projectTypes []domain.ProjectType) []projectTypeResponse {
	projectTypeResponses := []projectTypeResponse{}
	for _, projectType := range projectTypes {
		projectTypeResponses = append(projectTypeResponses, mapProjectTypeDomainToProjectTypeResponse(projectType))
	}
	return projectTypeResponses
}

func mapProjectTypeArrayDomainToSearchProjectTypesResponse(users []domain.ProjectType) searchProjectTypesResponse {
	return searchProjectTypesResponse{
		ProjectTypes: mapProjectTypeArrayDomainToProjectTypeArrayResponse(users),
	}
}

func mapSearchProjectTypesRequestToProjectTypeFilterDomain(request searchProjectTypesRequest) domain.ProjectTypeFilter {
	return domain.ProjectTypeFilter{
		WorkspaceId: request.WorkspaceId,
	}
}
