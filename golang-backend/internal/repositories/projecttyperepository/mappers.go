package projecttyperepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func DomainToDTO(projectType domain.ProjectType) projectTypeDTO {
	return projectTypeDTO{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func DTOToDomain(projectType projectTypeDTO) domain.ProjectType {
	return domain.ProjectType{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func DTOArrayToDomainArray(projectTypes []projectTypeDTO) []domain.ProjectType {
	results := make([]domain.ProjectType, 0, len(projectTypes))
	for _, projectType := range projectTypes {
		results = append(results, DTOToDomain(projectType))
	}
	return results
}
