package taskrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type TaskRepository interface {
	GetNextTaskIndex(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) (uint16, error)
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, task domain.Task) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, task domain.Task) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, deletedAt int64) error
	UnDelete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error
	Purge(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) (domain.Task, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error)
}

type taskRepository struct{}

func NewRepository() TaskRepository {
	return &taskRepository{}
}

func (repository *taskRepository) GetNextTaskIndex(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) (uint16, error) {
	var taskIndex uint16
	err := dbExecutor.QueryRowContext(
		ctx,
		`
			UPDATE project_task_counter SET
				next_task_index = next_task_index + 1
			WHERE
				project_id = ?
			RETURNING next_task_index - 1;
        `,
		projectID,
	).Scan(&taskIndex)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, domain.NotFoundError
		}
		return 0, err
	}
	return taskIndex, nil
}

func (repository *taskRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, task domain.Task) error {
	dto := toDTO(task)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO tasks
				(id, project_id, task_index, summary, description, creator_id, created_at, updated_at, deleted_at, started_at, finished_at, due_at, estimated_time, priority_id, status_id, cover_attachment_id)
			VALUES
				(?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, NULL)
        `,
		dto.ID,
		projectID,
		dto.Index,
		dto.Summary,
		dto.Description,
		dto.CreatorID,
		dto.CreatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.EstimatedTime,
		dto.PriorityID,
		dto.StatusID,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *taskRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, task domain.Task) error {
	dto := toDTO(task)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE tasks SET
				summary = ?,
				description = ?,
				updated_at = ?,
				started_at = ?,
				finished_at = ?,
				due_at = ?,
				estimated_time = ?,
				priority_id = ?,
				status_id = ?
			WHERE
				id = ?
			AND
				deleted_at IS NOT NULL
        `,
		dto.Summary,
		dto.Description,
		dto.UpdatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.EstimatedTime,
		dto.PriorityID,
		dto.StatusID,
		dto.ID,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *taskRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, deletedAt int64) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE tasks SET
				deleted_at = ?
			WHERE
				id = ?
        `,
		deletedAt,
		taskID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *taskRepository) UnDelete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE tasks SET
				deleted_at = NULL
			WHERE
				id = ?
        `,
		taskID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *taskRepository) Purge(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM tasks
			WHERE
				id = ?
        `,
		taskID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *taskRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) (domain.Task, error) {
	var dto taskDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                T.id,
				T.project_id,
				T.task_index,
				P.slug || "-" || T.task_index AS task_slug,
				T.summary,
				T.description,
				T.created_at,
				T.updated_at,
				T.deleted_at,
				T.started_at,
				T.finished_at,
				T.due_at,
				T.estimated_time,
				T.status_id,
				TS.name AS status_name,
				TS.item_hex_color AS status_hex_color,
				T.priority_id,
				TP.name AS priority_name,
				TP.item_hex_color AS priority_hex_color,
				T.creator_id,
				U.name AS creator_name,
				IFNULL(TN.notes_count, 0) AS notes_count,
				IFNULL(TA.attachments_count, 0) AS attachments_count,
				IFNULL(THO.history_operations_count, 0) AS history_operations_count,
				IFNULL(TTT.time_trackings_count, 0) AS time_trackings_count,
				IFNULL(TTT_SUM.total_spent_time, 0) AS total_spent_time
            FROM tasks T
			INNER JOIN projects P on P.id = T.project_id
			INNER JOIN task_priorities TP ON TP.id = T.priority_id
			INNER JOIN task_statuses TS ON TS.id = T.status_id
			INNER JOIN users U ON U.ID = T.creator_id
			LEFT JOIN (
				SELECT task_id, COUNT(*) AS time_trackings_count
				FROM task_time_trackings
				GROUP BY task_id
			) TTT ON TTT.task_id = T.id
			LEFT JOIN (
    			SELECT task_id, COUNT(*) AS notes_count
    			FROM task_notes
    			GROUP BY task_id
			) TN ON TN.task_id = T.id
			LEFT JOIN (
				SELECT task_id, COUNT(*) AS attachments_count
				FROM task_attachments
				GROUP BY task_id
			) TA ON TA.task_id = T.id
			LEFT JOIN (
				SELECT task_id, COUNT(*) as history_operations_count
				FROM history_operations
				GROUP BY task_id
			) THO ON THO.task_id = T.id
			LEFT JOIN (
				SELECT task_id, SUM(spent_time) as total_spent_time
				FROM task_time_trackings
				GROUP BY task_id
			) TTT_SUM ON TTT_SUM.task_id = T.id
            WHERE
				T.id = ?
        `,
		taskID).Scan(
		&dto.ID,
		&dto.projectID,
		&dto.Index,
		&dto.Slug,
		&dto.Summary,
		&dto.Description,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.DeletedAt,
		&dto.StartedAt,
		&dto.FinishedAt,
		&dto.DueAt,
		&dto.EstimatedTime,
		&dto.StatusID,
		&dto.StatusName,
		&dto.StatusHexColor,
		&dto.PriorityID,
		&dto.PriorityName,
		&dto.PriorityHexColor,
		&dto.CreatorID,
		&dto.CreatorName,
		&dto.NotesCount,
		&dto.AttachmentsCount,
		&dto.HistoryOperationsCount,
		&dto.TimeTrackingsCount,
		&dto.TotalSpentTime,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.NotFoundError
		}
		return domain.Task{}, err
	}
	return toDomain(dto), err
}

func (repository *taskRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
			T.id,
			T.project_id,
			T.task_index,
			P.slug || "-" || T.task_index AS task_slug,
			T.summary,
			T.description,
			T.created_at,
			T.updated_at,
			T.deleted_at,
			T.started_at,
			T.finished_at,
			T.due_at,
			T.estimated_time,
			T.status_id,
			TS.name AS status_name,
			TS.item_hex_color AS status_hex_color,
			T.priority_id,
			TP.name AS priority_name,
			TP.item_hex_color AS priority_hex_color,
			T.creator_id,
			U.name AS creator_name
		FROM tasks T
	`
	sqlQueryInnerJoins := `
		INNER JOIN projects P on P.id = T.project_id
		INNER JOIN task_priorities TP ON TP.id = T.priority_id
		INNER JOIN task_statuses TS ON TS.id = T.status_id
		INNER JOIN users U ON U.ID = T.creator_id
	`
	var field string
	switch order.Field {
	case "priority":
		field = "TP.name COLLATE NOCASE"
	case "status":
		field = "TS.item_index"
	case "summary":
		field = "P.summary COLLATE NOCASE"
	case "createdAt":
		field = "P.created_at"
	case "updatedAt":
		field = "P.updated_at"
	case "deletedAt":
		field = "P.deleted_at"
	case "startedAt":
		field = "P.started_at"
	case "finishedAt":
		field = "P.finished_at"
	case "dueAt":
		field = "P.due_at"
	case "createdBy":
		field = "U.name COLLATE NOCASE"
	default:
		field = "T.task_index"
	}
	var sort string
	switch order.Sort {
	case "DESC":
		sort = "DESC"
	case "ASC":
		sort = "ASC"
	default:
		sort = "ASC"
	}
	sqlOrder := fmt.Sprintf(" ORDER BY %s %s ", field, sort)
	sqlWhere := ""
	var sqlWhereConditions []string
	if filterDTO.ProjectID != nil && len(*filterDTO.ProjectID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.project_id = ?")
		filterArgs = append(filterArgs, *filterDTO.ProjectID)
	}
	if filterDTO.Summary != nil && len(*filterDTO.Summary) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.summary LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Summary+"%")
	}
	if filterDTO.PriorityID != nil && len(*filterDTO.PriorityID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.priority_id = ?")
		filterArgs = append(filterArgs, *filterDTO.PriorityID)
	}
	if filterDTO.StatusID != nil && len(*filterDTO.StatusID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.status_id = ?")
		filterArgs = append(filterArgs, *filterDTO.StatusID)
	}
	if filterDTO.CreatedAt != nil {
		if filterDTO.CreatedAt.From != nil && *filterDTO.CreatedAt.From > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "T.created_at >= ?")
			filterArgs = append(filterArgs, filterDTO.CreatedAt.From)
		}
		if filterDTO.CreatedAt.To != nil && *filterDTO.CreatedAt.To > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "T.created_at <= ?")
			filterArgs = append(filterArgs, filterDTO.CreatedAt.To)
		}
	}
	// TODO: updatedat, deletedat, startedat, finishedat, dueat
	if filterDTO.CreatedByUserID != nil && len(*filterDTO.CreatedByUserID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.creator_id = ?")
		filterArgs = append(filterArgs, *filterDTO.CreatedByUserID)
	}
	if filterDTO.ViewByUserID != nil && len(*filterDTO.ViewByUserID) > 0 {
		// TODO: onlw show tasks with user view permission
	}
	if len(sqlWhereConditions) > 0 {
		sqlWhere = " WHERE " + strings.Join(sqlWhereConditions, " AND ")
	}
	queryArgs = append(queryArgs, filterArgs...)
	var sqlLimit string
	if pager.Enabled() {
		sqlLimit = " LIMIT ? OFFSET ? "
		queryArgs = append(queryArgs, pager.Limit(), pager.Offset())
	} else {
		sqlLimit = ""
	}
	sqlQuery = fmt.Sprintf("%s %s %s %s %s ", sqlQuery, sqlQueryInnerJoins, sqlWhere, sqlOrder, sqlLimit)
	rows, err := dbExecutor.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	dtos := make([]taskDTO, 0)
	for rows.Next() {
		var dto taskDTO
		if err := rows.Scan(
			&dto.ID,
			&dto.projectID,
			&dto.Index,
			&dto.Slug,
			&dto.Summary,
			&dto.Description,
			&dto.CreatedAt,
			&dto.UpdatedAt,
			&dto.DeletedAt,
			&dto.StartedAt,
			&dto.FinishedAt,
			&dto.DueAt,
			&dto.EstimatedTime,
			&dto.StatusID,
			&dto.StatusName,
			&dto.StatusHexColor,
			&dto.PriorityID,
			&dto.PriorityName,
			&dto.PriorityHexColor,
			&dto.CreatorID,
			&dto.CreatorName,
		); err != nil {
			return nil, browser.Result{}, err
		}
		dtos = append(dtos, dto)
	}
	var totalResults int

	if pager.Enabled() {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS total_tasks
			FROM tasks T
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = dbExecutor.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(dtos)
	}

	return toDomainArray(dtos), browser.NewResult(pager, totalResults), nil
}
