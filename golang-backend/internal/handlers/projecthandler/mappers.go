package projecthandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddProjectRequestToProjectDomain(request addProjectRequest) domain.Project {
	return domain.Project{
		ID:      request.ID,
		Key:     request.Key,
		Summary: request.Summary,
	}
}

func mapUpdateProjectRequestToProjectDomain(request updateProjectRequest) domain.Project {
	return domain.Project{
		ID:      request.ID,
		Key:     request.Key,
		Summary: request.Summary,
	}
}

func mapProjectDomainToProjectResponse(project domain.Project) projectResponse {
	return projectResponse{
		ID:      project.ID,
		Key:     project.Key,
		Summary: project.Summary,
	}
}

func mapProjectDomainToAddProjectResponse(project domain.Project) addProjectResponse {
	return addProjectResponse{
		Project: mapProjectDomainToProjectResponse(project),
	}
}

func mapProjectDomainToUpdateProjectResponse(project domain.Project) updateProjectResponse {
	return updateProjectResponse{
		Project: mapProjectDomainToProjectResponse(project),
	}
}

func mapProjectDomainToGetProjectResponse(project domain.Project) getProjectResponse {
	return getProjectResponse{
		Project: mapProjectDomainToProjectResponse(project),
	}
}

func mapProjectArrayDomainToProjectArrayResponse(projects []domain.Project) []projectResponse {
	var projectResponses []projectResponse
	for _, project := range projects {
		projectResponses = append(projectResponses, mapProjectDomainToProjectResponse(project))
	}
	return projectResponses
}

func mapProjectArrayDomainToSearchProjectsResponse(users []domain.Project) searchProjectsResponse {
	return searchProjectsResponse{
		Projects: mapProjectArrayDomainToProjectArrayResponse(users),
	}
}
