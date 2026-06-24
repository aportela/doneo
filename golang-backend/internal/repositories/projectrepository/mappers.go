package projectrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(project domain.Project) projectDTO {
	return projectDTO{
		ID:                     project.ID,
		Slug:                   project.Slug,
		Summary:                project.Summary,
		Description:            utils.StrPtrToSQLNullStr(project.Description),
		CreatorID:              project.CreatedBy.ID,
		CreatorName:            project.CreatedBy.Name,
		CreatedAt:              project.CreatedAt.UnixMilli(),
		UpdatedAt:              utils.TimePtrToSQLNullInt64(project.UpdatedAt),
		DeletedAt:              utils.TimePtrToSQLNullInt64(project.DeletedAt),
		StartedAt:              utils.TimePtrToSQLNullInt64(project.StartedAt),
		FinishedAt:             utils.TimePtrToSQLNullInt64(project.FinishedAt),
		DueAt:                  utils.TimePtrToSQLNullInt64(project.DueAt),
		TypeID:                 project.Type.ID,
		TypeName:               project.Type.Name,
		TypeHexColor:           project.Type.HexColor,
		PriorityID:             project.Priority.ID,
		PriorityName:           project.Priority.Name,
		StatusID:               project.Status.ID,
		StatusName:             project.Status.Name,
		TasksCount:             project.TasksCount,
		PermissionsCount:       project.PermissionsCount,
		AttachmentsCount:       project.AttachmentsCount,
		NotesCount:             project.NotesCount,
		HistoryOperationsCount: project.HistoryOperationsCount,
	}
}

func toDomain(project projectDTO) domain.Project {
	return domain.Project{
		ID:                     project.ID,
		Slug:                   project.Slug,
		Summary:                project.Summary,
		Description:            utils.SQLStrPtr(project.Description),
		CreatedBy:              domain.UserBase{ID: project.CreatorID, Name: project.CreatorName},
		CreatedAt:              time.UnixMilli(project.CreatedAt),
		UpdatedAt:              utils.SQLNullInt64ToTimePtr(project.UpdatedAt),
		DeletedAt:              utils.SQLNullInt64ToTimePtr(project.DeletedAt),
		StartedAt:              utils.SQLNullInt64ToTimePtr(project.StartedAt),
		FinishedAt:             utils.SQLNullInt64ToTimePtr(project.FinishedAt),
		DueAt:                  utils.SQLNullInt64ToTimePtr(project.DueAt),
		Type:                   domain.ProjectType{ID: project.TypeID, Name: project.TypeName, HexColor: project.TypeHexColor},
		Priority:               domain.ProjectPriority{ID: project.PriorityID, Name: project.PriorityName, HexColor: project.PriorityHexColor},
		Status:                 domain.ProjectStatus{ID: project.StatusID, Name: project.StatusName, HexColor: project.StatusHexColor},
		TasksCount:             project.TasksCount,
		PermissionsCount:       project.PermissionsCount,
		AttachmentsCount:       project.AttachmentsCount,
		NotesCount:             project.NotesCount,
		HistoryOperationsCount: project.HistoryOperationsCount,
		PermissionsBitMask:     domain.Bitmask(project.PermissionsBitmask),
	}
}

func toDomainArray(projects []projectDTO) []domain.Project {
	results := []domain.Project{}
	for _, project := range projects {
		results = append(results, toDomain(project))
	}
	return results
}

func toFilterDTO(filter domain.SearchProjectFilter) searchFilterDTO {
	return searchFilterDTO{
		Slug:            filter.Slug,
		Summary:         filter.Summary,
		TypeID:          filter.TypeID,
		PriorityID:      filter.PriorityID,
		StatusID:        filter.StatusID,
		CreatedAt:       repositories.TimestampFilterToDTO(filter.CreatedAt),
		CreatedByUserID: filter.CreatedByUserID,
	}
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		if strings.Contains(sqlErr.Error(), "projects.slug") {
			return &domain.AlreadyExistsError{Field: "slug"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(slug)") {
			return &domain.ValidationError{Field: "slug"}
		} else if strings.Contains(sqlErr.Error(), "length(summary)") {
			return &domain.ValidationError{Field: "summary"}
		}
	}
	return err
}
