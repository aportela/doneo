package projectrepository

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
)

func DomainToDTO(project domain.Project) projectDTO {
	return projectDTO{
		ID:                     project.ID,
		Key:                    project.Key,
		Summary:                project.Summary,
		Description:            project.Description,
		CreatorId:              project.CreatedBy.ID,
		CreatorName:            project.CreatedBy.Name,
		CreatedAt:              project.CreatedAt,
		UpdatedAt:              project.UpdatedAt,
		StartedAt:              project.StartedAt,
		FinishedAt:             project.FinishedAt,
		DueAt:                  project.DueAt,
		TypeId:                 project.Type.ID,
		TypeName:               project.Type.Name,
		TypeHexColor:           project.Type.HexColor,
		PriorityId:             project.Priority.ID,
		PriorityName:           project.Priority.Name,
		StatusId:               project.Status.ID,
		StatusName:             project.Status.Name,
		TasksCount:             project.TasksCount,
		PermissionsCount:       project.PermissionsCount,
		AttachmentsCount:       project.AttachmentsCount,
		NotesCount:             project.NotesCount,
		HistoryOperationsCount: project.HistoryOperationsCount,
	}
}

func DTOToDomain(project projectDTO) domain.Project {
	return domain.Project{
		ID:                     project.ID,
		Key:                    project.Key,
		Summary:                project.Summary,
		Description:            project.Description,
		CreatedBy:              domain.UserBase{ID: project.CreatorId, Name: project.CreatorName},
		CreatedAt:              project.CreatedAt,
		UpdatedAt:              project.UpdatedAt,
		StartedAt:              project.StartedAt,
		FinishedAt:             project.FinishedAt,
		DueAt:                  project.DueAt,
		Type:                   domain.ProjectType{ID: project.TypeId, Name: project.TypeName, HexColor: project.TypeHexColor},
		Priority:               domain.ProjectPriority{ID: project.PriorityId, Name: project.PriorityName, HexColor: project.PriorityHexColor},
		Status:                 domain.ProjectStatus{ID: project.StatusId, Name: project.StatusName, HexColor: project.StatusHexColor},
		TasksCount:             project.TasksCount,
		PermissionsCount:       project.PermissionsCount,
		AttachmentsCount:       project.AttachmentsCount,
		NotesCount:             project.NotesCount,
		HistoryOperationsCount: project.HistoryOperationsCount,
	}
}

func DTOArrayToDomainArray(projects []projectDTO) []domain.Project {
	results := []domain.Project{}
	for _, project := range projects {
		results = append(results, DTOToDomain(project))
	}
	return results
}

func DomainFilterToDTO(filter domain.SearchProjectFilter) searchFilterDTO {

	return searchFilterDTO{
		Key:             filter.Key,
		Summary:         filter.Summary,
		TypeId:          filter.TypeId,
		PriorityId:      filter.PriorityId,
		StatusId:        filter.StatusId,
		CreatedAt:       repositories.TimestampFilterToDTO(filter.CreatedAt),
		CreatedByUserId: filter.CreatedByUserId,
	}
}
