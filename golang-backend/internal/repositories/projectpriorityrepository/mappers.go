package projectpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectPriorityDomainToProjectPriorityDTO(projectPriority domain.ProjectPriority) projectPriorityDTO {
	return projectPriorityDTO{
		ID:          projectPriority.ID,
		WorkspaceId: projectPriority.WorkspaceId,
		Name:        projectPriority.Name,
		Index:       projectPriority.Index,
		HexColor:    projectPriority.HexColor,
	}
}

func MapProjectPriorityDTOToProjectPriorityDomain(projectPriority projectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:          projectPriority.ID,
		WorkspaceId: projectPriority.WorkspaceId,
		Name:        projectPriority.Name,
		Index:       projectPriority.Index,
		HexColor:    projectPriority.HexColor,
	}
}

func MapProjectPriorityArrayDTOToProjectPriorityArrayDomain(projectPriority []projectPriorityDTO) []domain.ProjectPriority {
	results := []domain.ProjectPriority{}
	for _, projectPriority := range projectPriority {
		results = append(results, MapProjectPriorityDTOToProjectPriorityDomain(projectPriority))
	}
	return results
}

func MapProjectPriorityFilterDomainToProjectPriorityFilterDTO(filter domain.ProjectPriorityFilter) projectPriorityFilterDTO {
	return projectPriorityFilterDTO{
		WorkspaceId: filter.WorkspaceId,
	}
}

func MapProjectPriorityFilterDTOToProjectPriorityFilterDomain(filter projectPriorityFilterDTO) domain.ProjectPriorityFilter {
	return domain.ProjectPriorityFilter{
		WorkspaceId: filter.WorkspaceId,
	}
}
