package projectpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(projectPriority domain.ProjectPriority) projectPriorityDTO {
	return projectPriorityDTO{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func toDomain(projectPriority projectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func toDomainArray(projectPriorities []projectPriorityDTO) []domain.ProjectPriority {
	results := make([]domain.ProjectPriority, 0, len(projectPriorities))
	for _, projectPriority := range projectPriorities {
		results = append(results, toDomain(projectPriority))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectPrioritiesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
