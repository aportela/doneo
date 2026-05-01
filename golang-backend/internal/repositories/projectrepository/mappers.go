package projectrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapProjectDomainToProjectDTO(project domain.Project) projectDTO {
	return projectDTO{
		ID:           project.ID,
		WorkspaceId:  project.WorkspaceId,
		Key:          project.Key,
		Summary:      project.Summary,
		Description:  project.Description,
		CreatorId:    project.CreatedBy.ID,
		CreatorName:  project.CreatedBy.Name,
		CreatedAt:    project.CreatedAt,
		UpdatedAt:    project.UpdatedAt,
		StartedAt:    project.StartedAt,
		FinishedAt:   project.FinishedAt,
		DueAt:        project.DueAt,
		TypeId:       project.Type.ID,
		TypeName:     project.Type.Name,
		PriorityId:   project.Priority.ID,
		PriorityName: project.Priority.Name,
		StatusId:     project.Status.ID,
		StatusName:   project.Status.Name,
	}
}

func MapProjectDTOToProjectDomain(project projectDTO) domain.Project {
	return domain.Project{
		ID:          project.ID,
		WorkspaceId: project.WorkspaceId,
		Key:         project.Key,
		Summary:     project.Summary,
		Description: project.Description,
		CreatedBy:   domain.UserBase{ID: project.CreatorId, Name: project.CreatorName},
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
		StartedAt:   project.StartedAt,
		FinishedAt:  project.FinishedAt,
		DueAt:       project.DueAt,
		Type:        domain.ProjectType{ID: project.TypeId, Name: project.TypeName},
		Priority:    domain.ProjectPriority{ID: project.PriorityId, Name: project.PriorityName, Index: project.PriorityIndex, HexColor: project.PriorityHexColor},
		Status:      domain.ProjectStatus{ID: project.StatusId, Name: project.StatusName, Index: project.StatusIndex, HexColor: project.StatusHexColor},
	}
}

func MapProjectArrayDTOToProjectArrayDomain(projects []projectDTO) []domain.Project {
	var results []domain.Project
	for _, project := range projects {
		results = append(results, MapProjectDTOToProjectDomain(project))
	}
	return results
}
