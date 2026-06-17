package taskrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(task domain.Task) taskDTO {
	return taskDTO{
		ID:                     task.ID,
		projectID:              task.ProjectID,
		Index:                  task.Index,
		Slug:                   task.Slug,
		Summary:                task.Summary,
		Description:            utils.StrPtrToSQLNullStr(task.Description),
		CreatorID:              task.CreatedBy.ID,
		CreatorName:            task.CreatedBy.Name,
		CreatedAt:              task.CreatedAt.UnixMilli(),
		UpdatedAt:              utils.TimePtrToSQLNullInt64(task.UpdatedAt),
		DeletedAt:              utils.TimePtrToSQLNullInt64(task.DeletedAt),
		StartedAt:              utils.TimePtrToSQLNullInt64(task.StartedAt),
		FinishedAt:             utils.TimePtrToSQLNullInt64(task.FinishedAt),
		DueAt:                  utils.TimePtrToSQLNullInt64(task.DueAt),
		PriorityID:             task.Priority.ID,
		PriorityName:           task.Priority.Name,
		StatusID:               task.Status.ID,
		StatusName:             task.Status.Name,
		AttachmentsCount:       task.AttachmentsCount,
		NotesCount:             task.NotesCount,
		HistoryOperationsCount: task.HistoryOperationsCount,
	}
}

func toDomain(task taskDTO) domain.Task {
	return domain.Task{
		ID:                     task.ID,
		ProjectID:              task.projectID,
		Index:                  task.Index,
		Slug:                   task.Slug,
		Summary:                task.Summary,
		Description:            utils.SQLStrPtr(task.Description),
		CreatedBy:              domain.UserBase{ID: task.CreatorID, Name: task.CreatorName},
		CreatedAt:              time.UnixMilli(task.CreatedAt),
		UpdatedAt:              utils.SQLNullInt64ToTimePtr(task.UpdatedAt),
		DeletedAt:              utils.SQLNullInt64ToTimePtr(task.DeletedAt),
		StartedAt:              utils.SQLNullInt64ToTimePtr(task.StartedAt),
		FinishedAt:             utils.SQLNullInt64ToTimePtr(task.FinishedAt),
		DueAt:                  utils.SQLNullInt64ToTimePtr(task.DueAt),
		Priority:               domain.TaskPriority{ID: task.PriorityID, Name: task.PriorityName, HexColor: task.PriorityHexColor},
		Status:                 domain.TaskStatus{ID: task.StatusID, Name: task.StatusName, HexColor: task.StatusHexColor},
		AttachmentsCount:       task.AttachmentsCount,
		NotesCount:             task.NotesCount,
		HistoryOperationsCount: task.HistoryOperationsCount,
	}
}

func toDomainArray(tasks []taskDTO) []domain.Task {
	results := []domain.Task{}
	for _, project := range tasks {
		results = append(results, toDomain(project))
	}
	return results
}

func toFilterDTO(filter domain.SearchTaskFilter) searchFilterDTO {
	return searchFilterDTO{
		ProjectID:       filter.ProjectID,
		Summary:         filter.Summary,
		PriorityID:      filter.PriorityID,
		StatusID:        filter.StatusID,
		CreatedByUserID: filter.CreatedByUserID,
	}
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(summary)") {
			return &domain.ValidationError{Field: "summary"}
		}
	}
	return err
}
