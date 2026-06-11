package taskstatushandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func requestFlagsToDomainFlagsBitmask(flags statusFlags) domain.Bitmask {
	bitmaskFlag := domain.Bitmask(0)
	if flags.DefaultStatusOnCreation {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagDefaultOnCreate)
	}
	if flags.FillEmptyStartDate {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagFillEmptyStartDate)
	}
	if flags.SetStartDate {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagSetStartDate)
	}
	if flags.FillEmptyFinishDate {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagFillEmptyFinishDate)
	}
	if flags.SetFinishDate {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagSetFinishDate)
	}
	if flags.UnsetFinishDateOnLeave {
		bitmaskFlag.AddFlag(domain.TaskStatusFlagUnsetFinishDateOnLeave)
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
		DefaultStatusOnCreation: flag.HasFlag(domain.TaskStatusFlagDefaultOnCreate),
		FillEmptyStartDate:      flag.HasFlag(domain.TaskStatusFlagFillEmptyStartDate),
		SetStartDate:            flag.HasFlag(domain.TaskStatusFlagSetStartDate),
		FillEmptyFinishDate:     flag.HasFlag(domain.TaskStatusFlagFillEmptyFinishDate),
		SetFinishDate:           flag.HasFlag(domain.TaskStatusFlagSetFinishDate),
		UnsetFinishDateOnLeave:  flag.HasFlag(domain.TaskStatusFlagUnsetFinishDateOnLeave),
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
