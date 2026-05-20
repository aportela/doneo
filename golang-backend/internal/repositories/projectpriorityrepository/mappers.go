package projectpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func DomainToDTO(projectPriority domain.ProjectPriority) projectPriorityDTO {
	return projectPriorityDTO{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
		Index:    projectPriority.Index,
	}
}

func DTOToDomain(projectPriority projectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		HexColor: projectPriority.HexColor,
		Index:    projectPriority.Index,
	}
}

func DTOArrayToDomainArray(projectPriorities []projectPriorityDTO) []domain.ProjectPriority {
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
