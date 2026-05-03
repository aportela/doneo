package projectstatusrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectStatusDomainToProjectStatusDTO(projectStatus domain.ProjectStatus) projectStatusDTO {
	return projectStatusDTO{
		ID:          projectStatus.ID,
		WorkspaceId: projectStatus.WorkspaceId,
		Name:        projectStatus.Name,
		Index:       projectStatus.Index,
		HexColor:    projectStatus.HexColor,
	}
}

func MapProjectStatusDTOToProjectStatusDomain(projectStatus projectStatusDTO) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:          projectStatus.ID,
		WorkspaceId: projectStatus.WorkspaceId,
		Name:        projectStatus.Name,
		Index:       projectStatus.Index,
		HexColor:    projectStatus.HexColor,
	}
}

func MapProjectStatusArrayDTOToProjectStatusArrayDomain(projectStatuses []projectStatusDTO) []domain.ProjectStatus {
	results := []domain.ProjectStatus{}
	for _, projectStatus := range projectStatuses {
		results = append(results, MapProjectStatusDTOToProjectStatusDomain(projectStatus))
	}
	return results
}

func MapProjectStatusFilterDomainToProjectStatusFilterDTO(filter domain.ProjectStatusFilter) projectStatusFilterDTO {
	return projectStatusFilterDTO{
		WorkspaceId: filter.WorkspaceId,
	}
}

func MapProjectStatusFilterDTOToProjectStatusFilterDomain(filter projectStatusFilterDTO) domain.ProjectStatusFilter {
	return domain.ProjectStatusFilter{
		WorkspaceId: filter.WorkspaceId,
	}
}
