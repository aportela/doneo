package workspacehandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddWorkspaceRequestToWorkspaceDomain(request addWorkspaceRequest) domain.Workspace {
	return domain.Workspace{
		ID:   request.ID,
		Name: request.Name,
	}
}

func mapWorkspaceDomainToWorkspaceResponse(workspace domain.Workspace) workspaceResponse {
	return workspaceResponse{
		ID:   workspace.ID,
		Name: workspace.Name,
	}
}

func mapWorkspaceDomainToAddWorkspaceResponse(workspace domain.Workspace) addWorkspaceResponse {
	return addWorkspaceResponse{
		Workspace: mapWorkspaceDomainToWorkspaceResponse(workspace),
	}
}

func mapWorkspaceArrayDomainToWorkspaceArrayResponse(workspaces []domain.Workspace) []workspaceResponse {
	workspaceResponses := []workspaceResponse{}
	for _, workspace := range workspaces {
		workspaceResponses = append(workspaceResponses, mapWorkspaceDomainToWorkspaceResponse(workspace))
	}
	return workspaceResponses
}

func mapWorkspaceArrayDomainToSearchWorkspacesResponse(workspaces []domain.Workspace) searchWorkspacesResponse {
	return searchWorkspacesResponse{
		Workspaces: mapWorkspaceArrayDomainToWorkspaceArrayResponse(workspaces),
	}
}
