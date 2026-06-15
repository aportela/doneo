package tasktimeentryhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.TaskTimeEntry {
	return domain.TaskTimeEntry{
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskTimeEntry {
	return domain.TaskTimeEntry{
		ID:           request.Id,
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func domainToResponse(taskTimeEntry domain.TaskTimeEntry) TaskTimeEntryResponse {
	return TaskTimeEntryResponse{
		Id:           taskTimeEntry.ID,
		User:         userhandler.BaseDomainToBaseResponse(taskTimeEntry.CreatedBy),
		CreatedAt:    taskTimeEntry.CreatedAt.UnixMilli(),
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func domainArrayToResponseArray(taskTimeEntries []domain.TaskTimeEntry) []TaskTimeEntryResponse {
	taskTimeEntriesResponse := []TaskTimeEntryResponse{}
	for _, taskTimeEntry := range taskTimeEntries {
		taskTimeEntriesResponse = append(taskTimeEntriesResponse, domainToResponse(taskTimeEntry))
	}
	return taskTimeEntriesResponse
}

func toSearchResponse(taskTimeEntries []domain.TaskTimeEntry) searchResponse {
	return searchResponse{
		ProjectPermissions: domainArrayToResponseArray(taskTimeEntries),
	}
}
