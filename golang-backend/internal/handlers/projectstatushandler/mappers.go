package projectstatushandler

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
	if flags.UnsetFinishDateOnLeave {
		bitmaskFlag.AddFlag(domain.ProjectStatusFlagUnsetFinishDateOnLeave)
	}
	return bitmaskFlag
}

func addRequestToDomain(request addRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
		ID:       request.ID,
		Name:     request.Name,
		HexColor: request.HexColor,
		Index:    request.Index,
		Flags:    requestFlagsToDomainFlagsBitmask(request.Flags),
	}
}

func updateRequestToDomain(request updateRequest) domain.ProjectStatus {
	return domain.ProjectStatus{
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
		UnsetFinishDateOnLeave:  flag.HasFlag(domain.ProjectStatusFlagUnsetFinishDateOnLeave),
	}
}

func DomainToResponse(projectStatus domain.ProjectStatus) ProjectStatusResponse {
	return ProjectStatusResponse{
		ID:       projectStatus.ID,
		Name:     projectStatus.Name,
		HexColor: projectStatus.HexColor,
		Index:    projectStatus.Index,
		Flags:    flagDomainToResponseFlags(projectStatus.Flags),
	}
}

func domainArrayToResponseArray(projectStatuses []domain.ProjectStatus) []ProjectStatusResponse {
	projectStatusResponses := []ProjectStatusResponse{}
	for _, projectStatus := range projectStatuses {
		projectStatusResponses = append(projectStatusResponses, DomainToResponse(projectStatus))
	}
	return projectStatusResponses
}

func toSearchResponse(projectStatuses []domain.ProjectStatus, pager browser.Result) searchResponse {
	return searchResponse{
		ProjectStatuses: domainArrayToResponseArray(projectStatuses),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
