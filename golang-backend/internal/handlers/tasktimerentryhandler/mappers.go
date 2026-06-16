package tasktimeentryhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.TaskTimerEntry {
	return domain.TaskTimerEntry{
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskTimerEntry {
	return domain.TaskTimerEntry{
		ID:           request.Id,
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func domainToResponse(taskTimeEntry domain.TaskTimerEntry) TaskTimeEntryResponse {
	return TaskTimeEntryResponse{
		Id:           taskTimeEntry.ID,
		User:         userhandler.BaseDomainToBaseResponse(taskTimeEntry.CreatedBy),
		CreatedAt:    taskTimeEntry.CreatedAt.UnixMilli(),
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func domainArrayToResponseArray(taskTimeEntries []domain.TaskTimerEntry) []TaskTimeEntryResponse {
	taskTimeEntriesResponse := []TaskTimeEntryResponse{}
	for _, taskTimeEntry := range taskTimeEntries {
		taskTimeEntriesResponse = append(taskTimeEntriesResponse, domainToResponse(taskTimeEntry))
	}
	return taskTimeEntriesResponse
}

func toSearchResponse(taskTimeEntries []domain.TaskTimerEntry) searchResponse {
	return searchResponse{
		ProjectPermissions: domainArrayToResponseArray(taskTimeEntries),
	}
}
