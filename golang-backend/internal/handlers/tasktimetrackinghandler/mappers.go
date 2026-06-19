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

func domainToResponse(taskTimeTracking domain.TaskTimeTracking) TaskTimeTrakingResponse {
	return TaskTimeTrakingResponse{
		ID:           taskTimeTracking.ID,
		CreatedBy:    userhandler.BaseDomainToBaseResponse(taskTimeTracking.CreatedBy),
		CreatedAt:    taskTimeTracking.CreatedAt.UnixMilli(),
		Summary:      taskTimeTracking.Summary,
		TotalSeconds: taskTimeTracking.TotalSeconds,
	}
}

func domainArrayToResponseArray(taskTimeEntries []domain.TaskTimeTracking) []TaskTimeTrakingResponse {
	taskTimeEntriesResponse := []TaskTimeTrakingResponse{}
	for _, taskTimeTracking := range taskTimeEntries {
		taskTimeEntriesResponse = append(taskTimeEntriesResponse, domainToResponse(taskTimeTracking))
	}
	return taskTimeEntriesResponse
}

func toSearchResponse(taskTimeEntries []domain.TaskTimeTracking) searchResponse {
	return searchResponse{
		TimeTrackings: domainArrayToResponseArray(taskTimeEntries),
	}
}
