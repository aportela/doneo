package taskstatushandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func addRequestToDomain(request addRequest) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
	}
}

func DomainToResponse(taskStatus domain.TaskStatus) TaskStatusResponse {
	return TaskStatusResponse{
		ID:       taskStatus.ID,
		Name:     taskStatus.Name,
		HexColor: taskStatus.HexColor,
		Index:    taskStatus.Index,
	}
}

func domainArrayToResponseArray(taskStatuses []domain.TaskStatus) []TaskStatusResponse {
	projectStatusResponses := []TaskStatusResponse{}
	for _, projectStatus := range taskStatuses {
		projectStatusResponses = append(projectStatusResponses, DomainToResponse(projectStatus))
	}
	return projectStatusResponses
}

func toSearchResponse(taskStatuses []domain.TaskStatus, pager browser.Result) searchResponse {
	return searchResponse{
		TaskStatuses: domainArrayToResponseArray(taskStatuses),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
