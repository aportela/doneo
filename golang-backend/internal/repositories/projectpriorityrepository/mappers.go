package projectpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProyectPriorityDomainToProyectPriorityDTO(projectPriority domain.ProjectPriority) projectPriorityDTO {
	return projectPriorityDTO{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		Index:    projectPriority.Index,
		HexColor: projectPriority.HexColor,
	}
}

func MapProyectPriorityDTOToProyectPriorityDomain(projectPriority projectPriorityDTO) domain.ProjectPriority {
	return domain.ProjectPriority{
		ID:       projectPriority.ID,
		Name:     projectPriority.Name,
		Index:    projectPriority.Index,
		HexColor: projectPriority.HexColor,
	}
}

func MapProyectPriorityArrayDTOToProyectPriorityArrayDomain(projectPriority []projectPriorityDTO) []domain.ProjectPriority {
	var results []domain.ProjectPriority
	for _, projectPriority := range projectPriority {
		results = append(results, MapProyectPriorityDTOToProyectPriorityDomain(projectPriority))
	}
	return results
}
