package tasktimetrackinghandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.TaskTimeTracking {
	return domain.TaskTimeTracking{
		Summary:   request.Summary,
		SpentTime: request.SpentTime,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskTimeTracking {
	return domain.TaskTimeTracking{
		ID:        request.ID,
		Summary:   request.Summary,
		SpentTime: request.SpentTime,
	}
}

func domainToResponse(taskTimeTracking domain.TaskTimeTracking) taskTimeTrakingResponse {
	return taskTimeTrakingResponse{
		ID:        taskTimeTracking.ID,
		CreatedBy: userhandler.BaseDomainToBaseResponse(taskTimeTracking.CreatedBy),
		CreatedAt: taskTimeTracking.CreatedAt.UnixMilli(),
		Summary:   taskTimeTracking.Summary,
		SpentTime: taskTimeTracking.SpentTime,
	}
}

func domainArrayToResponseArray(taskTimeEntries []domain.TaskTimeTracking) []taskTimeTrakingResponse {
	taskTimeEntriesResponse := []taskTimeTrakingResponse{}
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
