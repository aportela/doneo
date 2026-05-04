package projecttyperepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectTypeDomainToProjectTypeDTO(projectType domain.ProjectType) projectTypeDTO {
	return projectTypeDTO{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func MapProjectTypeDTOToProjectTypeDomain(projectType projectTypeDTO) domain.ProjectType {
	return domain.ProjectType{
		ID:       projectType.ID,
		Name:     projectType.Name,
		HexColor: projectType.HexColor,
	}
}

func MapProjectTypeArrayDTOToProjectTypeArrayDomain(projectTypes []projectTypeDTO) []domain.ProjectType {
	results := []domain.ProjectType{}
	for _, projectType := range projectTypes {
		results = append(results, MapProjectTypeDTOToProjectTypeDomain(projectType))
	}
	return results
}

func MapProjectTypeFilterDomainToProjectTypeFilterDTO(filter domain.ProjectTypeFilter) projectTypeFilterDTO {
	return projectTypeFilterDTO{}
}

func MapProjectTypeFilterDTOToProjectTypeFilterDomain(filter projectTypeFilterDTO) domain.ProjectTypeFilter {
	return domain.ProjectTypeFilter{}
}
