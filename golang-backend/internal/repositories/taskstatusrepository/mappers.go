package taskstatusrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func toDTO(taskStatus domain.TaskStatus) taskStatusDTO {
	return taskStatusDTO{
		ID:       taskStatus.ID,
		Name:     taskStatus.Name,
		HexColor: taskStatus.HexColor,
	}
}

func toDomain(taskStatus taskStatusDTO) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       taskStatus.ID,
		Name:     taskStatus.Name,
		HexColor: taskStatus.HexColor,
	}
}

func toDomainArray(taskStatuses []taskStatusDTO) []domain.TaskStatus {
	results := make([]domain.TaskStatus, 0, len(taskStatuses))
	for _, projectStatus := range taskStatuses {
		results = append(results, toDomain(projectStatus))
	}
	return results
}

func toFilterDTO(filter domain.SearchTaskStatusesFilter) searchFilterDTO {
	return searchFilterDTO{
		Name: filter.Name,
	}
}
