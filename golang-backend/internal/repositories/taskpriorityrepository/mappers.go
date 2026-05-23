package taskpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func DomainToDTO(taskPriority domain.TaskPriority) taskPriorityDTO {
	return taskPriorityDTO{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
	}
}

func DTOToDomain(taskPriority taskPriorityDTO) domain.TaskPriority {
	return domain.TaskPriority{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
	}
}

func DTOArrayToDomainArray(taskPriorities []taskPriorityDTO) []domain.TaskPriority {
	results := make([]domain.TaskPriority, 0, len(taskPriorities))
	for _, projectPriority := range taskPriorities {
		results = append(results, DTOToDomain(projectPriority))
	}
	return results
}

func DomainFilterToDTO(filter domain.SearchProjectPrioritiesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
