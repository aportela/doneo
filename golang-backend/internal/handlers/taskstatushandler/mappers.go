package taskstatushandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func requestFlagsToDomainFlagsBitmask(flags statusFlags) domain.Bitmask {
	bitmaskFlag := domain.Bitmask(0)
	if flags.DefaultStatusOnCreation {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagDefaultOnCreate)
	}
	if flags.FillEmptyStartDate {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagFillEmptyStartDate)
	}
	if flags.SetStartDate {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagSetStartDate)
	}
	if flags.FillEmptyFinishDate {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagFillEmptyFinishDate)
	}
	if flags.SetFinishDate {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagSetFinishDate)
	}
	return bitmaskFlag
}

func addRequestToDomain(request addRequest) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
		Flags:    requestFlagsToDomainFlagsBitmask(request.Flags),
	}
}

func updateRequestToDomain(request updateRequest) domain.TaskStatus {
	return domain.TaskStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
		Flags:    requestFlagsToDomainFlagsBitmask(request.Flags),
	}
}

func flagDomainToResponseFlags(flag domain.Bitmask) statusFlags {
	return statusFlags{
		DefaultStatusOnCreation: flag.HasFlag(domain.ProjectStatusFlagDefaultOnCreate),
		FillEmptyStartDate:      flag.HasFlag(domain.ProjectStatusFlagFillEmptyStartDate),
		SetStartDate:            flag.HasFlag(domain.ProjectStatusFlagSetStartDate),
		FillEmptyFinishDate:     flag.HasFlag(domain.ProjectStatusFlagFillEmptyFinishDate),
		SetFinishDate:           flag.HasFlag(domain.ProjectStatusFlagSetFinishDate),
	}
}

func DomainToResponse(taskStatus domain.TaskStatus) TaskStatusResponse {
	return TaskStatusResponse{
		ID:       taskStatus.ID,
		Name:     taskStatus.Name,
		HexColor: taskStatus.HexColor,
		Index:    taskStatus.Index,
		Flags:    flagDomainToResponseFlags(taskStatus.Flags),
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
