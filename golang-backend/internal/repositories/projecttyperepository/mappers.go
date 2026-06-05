package projecttyperepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(projectType domain.ProjectType) ProjectTypeDTO {
	return ProjectTypeDTO{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func toDomain(projectType ProjectTypeDTO) domain.ProjectType {
	return domain.ProjectType{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func toDomainArray(projectTypes []ProjectTypeDTO) []domain.ProjectType {
	results := make([]domain.ProjectType, 0, len(projectTypes))
	for _, projectType := range projectTypes {
		results = append(results, toDomain(projectType))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectTypesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
