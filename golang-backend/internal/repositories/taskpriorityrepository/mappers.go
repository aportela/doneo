package taskpriorityrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(taskPriority domain.TaskPriority) taskPriorityDTO {
	return taskPriorityDTO{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
		Index:    taskPriority.Index,
	}
}

func toDomain(taskPriority taskPriorityDTO) domain.TaskPriority {
	return domain.TaskPriority{
		ID:       taskPriority.ID,
		Name:     taskPriority.Name,
		HexColor: taskPriority.HexColor,
		Index:    taskPriority.Index,
	}
}

func toDomainArray(taskPriorities []taskPriorityDTO) []domain.TaskPriority {
	results := make([]domain.TaskPriority, 0, len(taskPriorities))
	for _, projectPriority := range taskPriorities {
		results = append(results, toDomain(projectPriority))
	}
	return results
}

func toFilterDTO(filter domain.SearchTaskPrioritiesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
