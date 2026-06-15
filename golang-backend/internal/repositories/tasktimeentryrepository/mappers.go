package tasktimeentryrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
)

func toDTO(taskTimeEntry domain.TaskTimeEntry) taskTimeEntryDTO {
	return taskTimeEntryDTO{
		ID:           taskTimeEntry.ID,
		CreatedAt:    taskTimeEntry.CreatedAt.UnixMilli(),
		CreatorId:    taskTimeEntry.CreatedBy.ID,
		CreatorName:  taskTimeEntry.CreatedBy.Name,
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func toDomain(taskTimeEntry taskTimeEntryDTO) domain.TaskTimeEntry {
	return domain.TaskTimeEntry{
		ID:           taskTimeEntry.ID,
		CreatedAt:    time.UnixMilli(taskTimeEntry.CreatedAt),
		CreatedBy:    domain.UserBase{ID: taskTimeEntry.CreatorId, Name: taskTimeEntry.CreatorName},
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func toDomainArray(timers []taskTimeEntryDTO) []domain.TaskTimeEntry {
	results := make([]domain.TaskTimeEntry, 0, len(timers))
	for _, timer := range timers {
		results = append(results, toDomain(timer))
	}
	return results
}
