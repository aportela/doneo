package projectstatusrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectStatusDomainToProjectStatusDTO(projectStatus domain.ProjectStatus) projectStatusDTO {
	return projectStatusDTO{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		Index:    projectStatus.Index,
		HexColor: projectStatus.HexColor,
	}
}

func MapProjectStatusDTOToProjectStatusDomain(projectStatus projectStatusDTO) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		Index:    projectStatus.Index,
		HexColor: projectStatus.HexColor,
	}
}

func MapProjectStatusArrayDTOToProjectStatusArrayDomain(projectStatuses []projectStatusDTO) []domain.ProjectStatus {
	var results []domain.ProjectStatus
	for _, projectStatus := range projectStatuses {
		results = append(results, MapProjectStatusDTOToProjectStatusDomain(projectStatus))
	}
	return results
}
