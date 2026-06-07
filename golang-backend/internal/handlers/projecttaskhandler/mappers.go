package projecttaskhandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToDomain(request addRequest) domain.Task {
	return domain.Task{
		ID:          request.ID,
		Summary:     request.Summary,
		Description: request.Description,
		CreatedBy:   domain.UserBase{},
		Priority: domain.TaskPriority{
			ID: request.Priority.ID,
		},
		Status: domain.TaskStatus{
			ID: request.Status.ID,
		},
	}
}

func updateRequestToDomain(request updateRequest) domain.Task {
	return domain.Task{
		ID:          request.ID,
		Summary:     request.Summary,
		Description: request.Description,
		Priority: domain.TaskPriority{
			ID: request.Priority.ID,
		},
		Status: domain.TaskStatus{
			ID: request.Status.ID,
		},
		StartedAt:  utils.Int64PtrToTimePtr(request.StartedAt),
		FinishedAt: utils.Int64PtrToTimePtr(request.FinishedAt),
		DueAt:      utils.Int64PtrToTimePtr(request.DueAt),
	}
}

func DomainToResponse(task domain.Task) taskResponse {
	return taskResponse{
		ID:                     task.ID,
		Slug:                   task.Slug,
		Summary:                task.Summary,
		Description:            task.Description,
		CreatedBy:              userhandler.BaseDomainToBaseResponse(task.CreatedBy),
		CreatedAt:              task.CreatedAt.UnixMilli(),
		UpdatedAt:              utils.TimePtrToInt64Ptr(task.UpdatedAt),
		DeletedAt:              utils.TimePtrToInt64Ptr(task.DeletedAt),
		StartedAt:              utils.TimePtrToInt64Ptr(task.StartedAt),
		FinishedAt:             utils.TimePtrToInt64Ptr(task.FinishedAt),
		DueAt:                  utils.TimePtrToInt64Ptr(task.DueAt),
		Priority:               taskpriorityhandler.DomainToResponse(task.Priority),
		Status:                 taskstatushandler.DomainToResponse(task.Status),
		PermissionsCount:       task.PermissionsCount,
		AttachmentsCount:       task.AttachmentsCount,
		NotesCount:             task.NotesCount,
		HistoryOperationsCount: task.HistoryOperationsCount,
	}
}

func domainArrayToResponseArray(tasks []domain.Task) []taskResponse {
	taskResponses := []taskResponse{}
	for _, task := range tasks {
		taskResponses = append(taskResponses, DomainToResponse(task))
	}
	return taskResponses
}

func toSearchResponse(tasks []domain.Task, pager browser.Result) searchTasksResponse {
	return searchTasksResponse{
		Tasks: domainArrayToResponseArray(tasks),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
