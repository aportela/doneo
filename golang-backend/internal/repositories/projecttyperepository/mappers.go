package projecttyperepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectTypeDomainToProjectTypeDTO(projectType domain.ProjectType) projectTypeDTO {
	return projectTypeDTO{
		ID:   projectType.ID,
		Name: projectType.Name,
	}
}

func MapProjectTypeDTOToProjectTypeDomain(projectType projectTypeDTO) domain.ProjectType {
	return domain.ProjectType{
		ID:   projectType.ID,
		Name: projectType.Name,
	}
}

func MapProjectTypeArrayDTOToProjectTypeArrayDomain(projectTypes []projectTypeDTO) []domain.ProjectType {
	var results []domain.ProjectType
	for _, projectType := range projectTypes {
		results = append(results, MapProjectTypeDTOToProjectTypeDomain(projectType))
	}
	return results
}
