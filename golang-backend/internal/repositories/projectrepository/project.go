package projectrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectRepository interface {
	Add(ctx context.Context, project projectDTO) error
	Update(ctx context.Context, project projectDTO) error
	Get(ctx context.Context, id string) (projectDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectDTO, browser.Result, error)
}

type projectRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectRepository {
	return &projectRepository{database: database}
}

func (repository *projectRepository) Add(ctx context.Context, project projectDTO) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(
		ctx,
		`
            INSERT INTO projects
				(id, key, summary, description, creator_id, created_at, updated_at, started_at, finished_at, due_at, priority_id, status_id, type_id)
			VALUES
				(?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?)
        `,
		project.ID,
		project.Key,
		project.Summary,
		project.Description,
		project.CreatorId,
		project.CreatedAt,
		project.StartedAt,
		project.FinishedAt,
		project.DueAt,
		project.PriorityId,
		project.StatusId,
		project.TypeId,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "projects.key") {
				return &domain.AlreadyExistsError{Field: "key"}
			} else if strings.Contains(sqlErr.Error(), "projects.id") {
				return &domain.AlreadyExistsError{Field: "id"}
			}
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(key)") {
				return &domain.ValidationError{Field: "key"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
		}
	}
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		project.ID,
		domain.EventProjectCreated,
		project.CreatorId,
		project.CreatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *projectRepository) Update(ctx context.Context, project projectDTO) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	_, err = tx.ExecContext(
		ctx,
		`
            UPDATE projects SET
				key = ?,
				summary = ?,
				description = ?,
				updated_at = ?,
				started_at = ?,
				finished_at = ?,
				due_at = ?,
				priority_id = ?,
				status_id = ?,
				type_id = ?
			WHERE id = ?
        `,
		project.Key,
		project.Summary,
		project.Description,
		// TODO: use received value
		utils.CurrentMSTimestamp(),
		project.StartedAt,
		project.FinishedAt,
		project.DueAt,
		project.PriorityId,
		project.StatusId,
		project.TypeId,
		project.ID,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "projects.key") {
				return &domain.AlreadyExistsError{Field: "key"}
			} else if strings.Contains(sqlErr.Error(), "projects.id") {
				return &domain.AlreadyExistsError{Field: "id"}
			}
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(key)") {
				return &domain.ValidationError{Field: "key"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
		}
	}
	userId, _ := middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		project.ID,
		domain.EventProjectUpdated,
		userId,
		project.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *projectRepository) Delete(ctx context.Context, id string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM projects
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (repository *projectRepository) Get(ctx context.Context, id string) (projectDTO, error) {
	var project projectDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id,
				P.key,
				P.summary,
				P.description,
				P.created_at,
				P.updated_at,
				P.started_at,
				P.finished_at,
				P.due_at,
				P.status_id,
				PS.name AS status_name,
				PS.item_hex_color AS status_hex_color,
				P.priority_id,
				PP.name AS priority_name,
				PP.item_hex_color AS priority_hex_color,
				P.type_id,
				PT.name AS type_name,
				PT.item_hex_color AS type_hex_color,
				P.creator_id,
				U.name AS creator_name,
				IFNULL(PUR.permissions_count, 0) AS permissions_count,
				IFNULL(PN.notes_count, 0) AS notes_count,
				IFNULL(PA.attachments_count, 0) AS attachments_count,
				IFNULL(PHO.history_operations_count, 0) AS history_operations_count,
				0 AS tasks_count
            FROM projects P
			INNER JOIN project_priorities PP ON PP.id = P.priority_id
			INNER JOIN project_statuses PS ON PS.id = P.status_id
			INNER JOIN project_types PT ON PT.id = P.type_id
			INNER JOIN users U ON U.ID = P.creator_id
			LEFT JOIN (
    			SELECT project_id, COUNT(*) AS permissions_count
    			FROM project_user_role
    			GROUP BY project_id
			) PUR ON PUR.project_id = P.id
			LEFT JOIN (
    			SELECT project_id, COUNT(*) AS notes_count
    			FROM project_notes
    			GROUP BY project_id
			) PN ON PN.project_id = P.id
			 LEFT JOIN (
    			SELECT project_id, COUNT(*) AS attachments_count
    			FROM project_attachments
    			GROUP BY project_id
			) PA ON PA.project_id = P.id
			LEFT JOIN (
				SELECT project_id, COUNT(*) as history_operations_count
				FROM project_history_operations
				GROUP BY project_id
			) PHO ON PHO.project_id = P.id
            WHERE P.id = ?
			GROUP BY P.id
        `,
		id).Scan(
		&project.ID,
		&project.Key,
		&project.Summary,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.StartedAt,
		&project.FinishedAt,
		&project.DueAt,
		&project.StatusId,
		&project.StatusName,
		&project.StatusHexColor,
		&project.PriorityId,
		&project.PriorityName,
		&project.PriorityHexColor,
		&project.TypeId,
		&project.TypeName,
		&project.TypeHexColor,
		&project.CreatorId,
		&project.CreatorName,
		&project.PermissionsCount,
		&project.NotesCount,
		&project.AttachmentsCount,
		&project.HistoryOperationsCount,
		&project.TasksCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectDTO{}, domain.NotFoundError
		}
		return projectDTO{}, err
	}
	return project, err
}

func (repository *projectRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectDTO, browser.Result, error) {
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
                P.id,
				P.key,
				P.summary,
				P.description,
				P.created_at,
				P.updated_at,
				P.started_at,
				P.finished_at,
				P.due_at,
				P.status_id,
				PS.name AS status_name,
				PS.item_hex_color AS status_hex_color,
				P.priority_id,
				PP.name AS priority_name,
				PP.item_hex_color AS priority_hex_color,
				P.type_id,
				PT.name AS type_name,
				PT.item_hex_color AS type_hex_color,
				P.creator_id,
				U.name AS creator_name
            FROM projects P
	`
	sqlQueryInnerJoins := `
		INNER JOIN project_priorities PP ON PP.id = P.priority_id
		INNER JOIN project_statuses PS ON PS.id = P.status_id
		INNER JOIN project_types PT ON PT.id = P.type_id
		INNER JOIN users U ON U.ID = P.creator_id
	`
	var field string
	switch order.Field {
	case "key":
		field = "P.key COLLATE NOCASE"
	case "type":
		field = "PT.name COLLATE NOCASE"
	case "priority":
		field = "PP.name COLLATE NOCASE"
	case "status":
		field = "PS.name COLLATE NOCASE"
	case "summary":
		field = "P.summary COLLATE NOCASE"
	case "createdAt":
		field = "P.created_at"
	case "createdBy":
		field = "U.name COLLATE NOCASE"
	default:
		field = "P.key COLLATE NOCASE"
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
	if filter.Key != nil && len(*filter.Key) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.key LIKE ?")
		filterArgs = append(filterArgs, "%"+*filter.Key+"%")
	}
	if filter.Summary != nil && len(*filter.Summary) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.summary LIKE ?")
		filterArgs = append(filterArgs, "%"+*filter.Summary+"%")
	}
	if filter.TypeId != nil && len(*filter.TypeId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.type_id = ?")
		filterArgs = append(filterArgs, *filter.TypeId)
	}
	if filter.PriorityId != nil && len(*filter.PriorityId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.priority_id = ?")
		filterArgs = append(filterArgs, *filter.PriorityId)
	}
	if filter.StatusId != nil && len(*filter.StatusId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.status_id = ?")
		filterArgs = append(filterArgs, *filter.StatusId)
	}
	if filter.CreatedAt != nil {
		if filter.CreatedAt.From != nil && *filter.CreatedAt.From > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "P.created_at >= ?")
			filterArgs = append(filterArgs, filter.CreatedAt.From)
		}
		if filter.CreatedAt.To != nil && *filter.CreatedAt.To > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "P.created_at <= ?")
			filterArgs = append(filterArgs, filter.CreatedAt.To)
		}
	}
	if filter.CreatedByUserId != nil && len(*filter.CreatedByUserId) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.creator_id = ?")
		filterArgs = append(filterArgs, *filter.CreatedByUserId)
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
	projects := make([]projectDTO, 0)
	for rows.Next() {
		var project projectDTO
		if err := rows.Scan(
			&project.ID,
			&project.Key,
			&project.Summary,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
			&project.StartedAt,
			&project.FinishedAt,
			&project.DueAt,
			&project.StatusId,
			&project.StatusName,
			&project.StatusHexColor,
			&project.PriorityId,
			&project.PriorityName,
			&project.PriorityHexColor,
			&project.TypeId,
			&project.TypeName,
			&project.TypeHexColor,
			&project.CreatorId,
			&project.CreatorName,
		); err != nil {
			return nil, browser.Result{}, err
		}
		projects = append(projects, project)
	}
	var totalResults int

	if pager.Enabled() {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS total_projects
			FROM projects P
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
		totalResults = len(projects)
	}

	return projects, browser.NewResult(pager, totalResults), nil
}
