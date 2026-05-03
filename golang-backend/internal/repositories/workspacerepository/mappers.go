package workspacerepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapWorkspaceDomainToWorkspaceDTO(workspace domain.Workspace) workspaceDTO {
	return workspaceDTO{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Description: workspace.Description,
		CreatorId:   workspace.CreatedBy.ID,
		CreatorName: workspace.CreatedBy.Name,
		CreatedAt:   workspace.CreatedAt,
		UpdatedAt:   workspace.UpdatedAt,
	}
}

func MapWorkspaceDTOToWorkspaceDomain(workspace workspaceDTO) domain.Workspace {
	return domain.Workspace{
		ID:          workspace.ID,
		Name:        workspace.Name,
		Description: workspace.Description,
		CreatedBy: domain.UserBase{
			ID:   workspace.CreatorId,
			Name: workspace.CreatorName,
		},
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: workspace.UpdatedAt,
	}
}

func MapWorkspaceArrayDTOToWorkspaceArrayDomain(workspaces []workspaceDTO) []domain.Workspace {
	results := []domain.Workspace{}
	for _, workspace := range workspaces {
		results = append(results, MapWorkspaceDTOToWorkspaceDomain(workspace))
	}
	return results
}
