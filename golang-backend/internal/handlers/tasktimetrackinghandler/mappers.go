package tasktimetrackinghandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.TaskTimeTracking {
	return domain.TaskTimeTracking{
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskTimeTracking {
	return domain.TaskTimeTracking{
		ID:           request.ID,
		Summary:      request.Summary,
		TotalSeconds: request.TotalSeconds,
	}
}

func domainToResponse(taskTimeEntry domain.TaskTimeTracking) TaskTimeTrakingResponse {
	return TaskTimeTrakingResponse{
		ID:           taskTimeEntry.ID,
		User:         userhandler.BaseDomainToBaseResponse(taskTimeEntry.CreatedBy),
		CreatedAt:    taskTimeEntry.CreatedAt.UnixMilli(),
		Summary:      taskTimeEntry.Summary,
		TotalSeconds: taskTimeEntry.TotalSeconds,
	}
}

func domainArrayToResponseArray(taskTimeEntries []domain.TaskTimeTracking) []TaskTimeTrakingResponse {
	taskTimeEntriesResponse := []TaskTimeTrakingResponse{}
	for _, taskTimeEntry := range taskTimeEntries {
		taskTimeEntriesResponse = append(taskTimeEntriesResponse, domainToResponse(taskTimeEntry))
	}
	return taskTimeEntriesResponse
}

func toSearchResponse(taskTimeEntries []domain.TaskTimeTracking) searchResponse {
	return searchResponse{
		TimeTrackings: domainArrayToResponseArray(taskTimeEntries),
	}
}
