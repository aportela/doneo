package projecttaskrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type TaskRepository interface {
	GetNextTaskIndex(ctx context.Context, projectId string) (uint, error)
	Add(ctx context.Context, projectId string, task domain.Task) error
	Update(ctx context.Context, task domain.Task) error
	Get(ctx context.Context, id string) (domain.Task, error)
	Delete(ctx context.Context, id string, deletedAt int64) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error)
}

type taskRepository struct {
	database database.Database
}

func NewRepository(database database.Database) TaskRepository {
	return &taskRepository{database: database}
}

func (repository *taskRepository) GetNextTaskIndex(ctx context.Context, projectId string) (uint, error) {
	var taskIndex uint
	err := repository.database.QueryRowContext(
		ctx,
		`
			UPDATE project_task_counter SET
				next_task_index = next_task_index + 1
			WHERE project_id = ?
			RETURNING next_task_index - 1;
        `,
		projectId,
	).Scan(&taskIndex)
	// TODO: no rows in result set ?
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return 0, err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return 0, &domain.ValidationError{Field: "project_id"}
			}
			return 0, err
		default:
			return 0, err
		}
	}
	return taskIndex, nil
}

func (repository *taskRepository) Add(ctx context.Context, projectId string, task domain.Task) error {
	dto := toDTO(task)
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO tasks
				(id, project_id, task_index, summary, description, creator_id, created_at, updated_at, deleted_at, started_at, finished_at, due_at, priority_id, status_id)
			VALUES
				(?, ?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		projectId,
		dto.Index,
		dto.Summary,
		dto.Description,
		dto.CreatorId,
		dto.CreatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.PriorityId,
		dto.StatusId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO: check unique
			return err
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskRepository) Update(ctx context.Context, task domain.Task) error {
	dto := toDTO(task)
	_, err := repository.database.ExecContext(
		ctx,
		`
            UPDATE tasks SET
				summary = ?,
				description = ?,
				updated_at = ?,
				started_at = ?,
				finished_at = ?,
				due_at = ?,
				priority_id = ?,
				status_id = ?
			WHERE id = ?
        `,
		dto.Summary,
		dto.Description,
		dto.UpdatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.PriorityId,
		dto.StatusId,
		dto.ID,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskRepository) Delete(ctx context.Context, id string, deletedAt int64) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            UPDATE tasks SET
				deleted_at = ?
			WHERE id = ?
        `,
		deletedAt,
		id,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *taskRepository) Get(ctx context.Context, id string) (domain.Task, error) {
	var dto taskDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                T.id,
				T.task_index,
				P.key || "-" || T.task_index AS task_slug,
				T.summary,
				T.description,
				T.created_at,
				T.updated_at,
				T.deleted_at,
				T.started_at,
				T.finished_at,
				T.due_at,
				T.status_id,
				TS.name AS status_name,
				TS.item_hex_color AS status_hex_color,
				T.priority_id,
				TP.name AS priority_name,
				TP.item_hex_color AS priority_hex_color,
				T.creator_id,
				U.name AS creator_name,
				0 AS notes_count,
				0 AS attachments_count,
				0 AS history_operations_count
            FROM tasks T
			INNER JOIN projects P on P.id = T.project_id
			INNER JOIN task_priorities TP ON TP.id = T.priority_id
			INNER JOIN task_statuses TS ON TS.id = T.status_id
			INNER JOIN users U ON U.ID = T.creator_id
            WHERE T.id = ?
        `,
		id).Scan(
		&dto.ID,
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
		&dto.StatusId,
		&dto.StatusName,
		&dto.StatusHexColor,
		&dto.PriorityId,
		&dto.PriorityName,
		&dto.PriorityHexColor,
		&dto.CreatorId,
		&dto.CreatorName,
		&dto.NotesCount,
		&dto.AttachmentsCount,
		&dto.HistoryOperationsCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.NotFoundError
		}
		return domain.Task{}, err
	}
	return toDomain(dto), err
}

func (repository *taskRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskFilter) ([]domain.Task, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
			T.id,
			T.task_index,
			P.key || "-" || T.task_index AS task_slug,
			T.summary,
			T.description,
			T.created_at,
			T.updated_at,
			T.deleted_at,
			T.started_at,
			T.finished_at,
			T.due_at,
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
		field = "TS.name COLLATE NOCASE"
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
	if filterDTO.Summary != nil && len(*filterDTO.Summary) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.summary LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Summary+"%")
	}
	if filterDTO.PriorityId != nil && len(*filterDTO.PriorityId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.priority_id = ?")
		filterArgs = append(filterArgs, *filterDTO.PriorityId)
	}
	if filterDTO.StatusId != nil && len(*filterDTO.StatusId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.status_id = ?")
		filterArgs = append(filterArgs, *filterDTO.StatusId)
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
	if filterDTO.CreatedByUserId != nil && len(*filterDTO.CreatedByUserId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "T.creator_id = ?")
		filterArgs = append(filterArgs, *filterDTO.CreatedByUserId)
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
	rows, err := repository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	dtos := make([]taskDTO, 0)
	for rows.Next() {
		var dto taskDTO
		if err := rows.Scan(
			&dto.ID,
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
			&dto.StatusId,
			&dto.StatusName,
			&dto.StatusHexColor,
			&dto.PriorityId,
			&dto.PriorityName,
			&dto.PriorityHexColor,
			&dto.CreatorId,
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
			FROM tasks P
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = repository.database.QueryRowContext(
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
