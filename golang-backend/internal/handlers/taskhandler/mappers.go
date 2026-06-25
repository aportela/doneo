package taskhandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
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
		EstimatedTime: request.EstimatedTime,
		Tags:          request.Tags,
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
		StartedAt:     utils.Int64PtrToTimePtr(request.StartedAt),
		FinishedAt:    utils.Int64PtrToTimePtr(request.FinishedAt),
		DueAt:         utils.Int64PtrToTimePtr(request.DueAt),
		EstimatedTime: request.EstimatedTime,
		Tags:          request.Tags,
	}
}

func patchRequestToDomain(request patchRequest) domain.Task {
	return domain.Task{
		ID: request.ID,
		Status: domain.TaskStatus{
			ID: request.Status.ID,
		},
	}
}

func DomainToResponse(task domain.Task) taskResponse {
	return taskResponse{
		ID:                     task.ID,
		ProjectID:              task.ProjectID,
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
		EstimatedTime:          task.EstimatedTime,
		TotalSpentTime:         task.TotalSpentTime,
		Priority:               taskpriorityhandler.DomainToResponse(task.Priority),
		Status:                 taskstatushandler.DomainToResponse(task.Status),
		Tags:                   task.Tags,
		PermissionsCount:       task.PermissionsCount,
		AttachmentsCount:       task.AttachmentsCount,
		NotesCount:             task.NotesCount,
		HistoryOperationsCount: task.HistoryOperationsCount,
		TimeTrackingsCount:     task.TimeTrackingsCount,
		AllowedOperations: projecthandler.ProjectAllowedOperationsResponse{
			AllowViewProject:   task.PermissionsBitMask.HasFlag(domain.PermissionViewProject),
			AllowUpdateProject: task.PermissionsBitMask.HasFlag(domain.PermissionUpdateProject),
			AllowDeleteProject: task.PermissionsBitMask.HasFlag(domain.PermissionDeleteProject),
			AllowAddTask:       task.PermissionsBitMask.HasFlag(domain.PermissionAddTask),
			AllowUpdateTask:    task.PermissionsBitMask.HasFlag(domain.PermissionUpdateTask),
			AllowDeleteTask:    task.PermissionsBitMask.HasFlag(domain.PermissionDeleteTask),
			AllowViewTask:      task.PermissionsBitMask.HasFlag(domain.PermissionViewTask),
		},
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
