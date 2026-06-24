package projecthandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToDomain(request addRequest) domain.Project {
	return domain.Project{
		ID:          request.ID,
		Slug:        request.Slug,
		Summary:     request.Summary,
		Description: request.Description,
		CreatedBy:   domain.UserBase{},
		Type: domain.ProjectType{
			ID: request.Type.ID,
		},
		Priority: domain.ProjectPriority{
			ID: request.Priority.ID,
		},
		Status: domain.ProjectStatus{
			ID: request.Status.ID,
		},
	}
}

func updateRequestToDomain(request updateRequest) domain.Project {
	return domain.Project{
		ID:          request.ID,
		Slug:        request.Slug,
		Summary:     request.Summary,
		Description: request.Description,
		Type: domain.ProjectType{
			ID: request.Type.ID,
		},
		Priority: domain.ProjectPriority{
			ID: request.Priority.ID,
		},
		Status: domain.ProjectStatus{
			ID: request.Status.ID,
		},
		StartedAt:  utils.Int64PtrToTimePtr(request.StartedAt),
		FinishedAt: utils.Int64PtrToTimePtr(request.FinishedAt),
		DueAt:      utils.Int64PtrToTimePtr(request.DueAt),
	}
}

func patchRequestToDomain(request patchRequest) domain.Project {
	return domain.Project{
		ID: request.ID,
		Status: domain.ProjectStatus{
			ID: request.Status.ID,
		},
	}
}

func DomainToResponse(project domain.Project) projectResponse {
	return projectResponse{
		ID:                     project.ID,
		SLug:                   project.Slug,
		Summary:                project.Summary,
		Description:            project.Description,
		CreatedBy:              userhandler.BaseDomainToBaseResponse(project.CreatedBy),
		CreatedAt:              project.CreatedAt.UnixMilli(),
		UpdatedAt:              utils.TimePtrToInt64Ptr(project.UpdatedAt),
		DeletedAt:              utils.TimePtrToInt64Ptr(project.DeletedAt),
		StartedAt:              utils.TimePtrToInt64Ptr(project.StartedAt),
		FinishedAt:             utils.TimePtrToInt64Ptr(project.FinishedAt),
		DueAt:                  utils.TimePtrToInt64Ptr(project.DueAt),
		Type:                   projecttypehandler.DomainToResponse(project.Type),
		Priority:               projectpriorityhandler.DomainToResponse(project.Priority),
		Status:                 projectstatushandler.DomainToResponse(project.Status),
		TasksCount:             project.TasksCount,
		PermissionsCount:       project.PermissionsCount,
		AttachmentsCount:       project.AttachmentsCount,
		NotesCount:             project.NotesCount,
		HistoryOperationsCount: project.HistoryOperationsCount,
		AllowedOperations: projectAllowedOperationsResponse{
			AllowViewProject:   project.PermissionsBitMask.HasFlag(domain.PermissionViewProject),
			AllowUpdateProject: project.PermissionsBitMask.HasFlag(domain.PermissionUpdateProject),
			AllowDeleteProject: project.PermissionsBitMask.HasFlag(domain.PermissionDeleteProject),
		},
	}
}

func domainArrayToResponseArray(projects []domain.Project) []projectResponse {
	projectResponses := []projectResponse{}
	for _, project := range projects {
		projectResponses = append(projectResponses, DomainToResponse(project))
	}
	return projectResponses
}

func toSearchResponse(users []domain.Project, pager browser.Result) searchProjectsResponse {
	return searchProjectsResponse{
		Projects: domainArrayToResponseArray(users),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
