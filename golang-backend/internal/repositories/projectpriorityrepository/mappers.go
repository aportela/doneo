package projectpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func DomainToDTO(projectPriority domain.ProjectPriority) ProjectPriorityDTO {
	return ProjectPriorityDTO{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func DTOToDomain(projectPriority ProjectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
	}
}

func DTOArrayToDomainArray(projectPriorities []ProjectPriorityDTO) []domain.ProjectPriority {
	results := make([]domain.ProjectPriority, 0, len(projectPriorities))
	for _, projectPriority := range projectPriorities {
		results = append(results, DTOToDomain(projectPriority))
	}
	return results
}

func DomainFilterToDTO(filter domain.SearchProjectPrioritiesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
