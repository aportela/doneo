package projectstatusrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(projectStatus domain.ProjectStatus) projectStatusDTO {
	return projectStatusDTO{
		ID:           projectStatus.ID,
		Name:         projectStatus.Name,
		HexColor:     projectStatus.HexColor,
		Index:        projectStatus.Index,
		FlagsBitmask: uint64(projectStatus.Flags),
	}
}

func toDomain(projectStatus projectStatusDTO) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		HexColor: projectStatus.HexColor,
		Index:    projectStatus.Index,
		Flags:    domain.Bitmask(projectStatus.FlagsBitmask),
	}
}

func toDomainArray(projectStatuses []projectStatusDTO) []domain.ProjectStatus {
	results := make([]domain.ProjectStatus, 0, len(projectStatuses))
	for _, projectStatus := range projectStatuses {
		results = append(results, toDomain(projectStatus))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectStatusesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
