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
)

type ProjectRepository interface {
	AddTaskCounter(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) error
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, project domain.Project) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, project domain.Project) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, deletedAt int64) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) (domain.Project, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.PagerQuery, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.PagerResult, error)
	GetCurrentProjects(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string, count uint) ([]domain.Project, error)
}

type projectRepository struct{}

func NewRepository() ProjectRepository {
	return &projectRepository{}
}

func (repository *projectRepository) AddTaskCounter(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
			INSERT INTO project_task_counter
				(project_id, next_task_index)
			VALUES
				(?, 1)
		`,
		projectID,
	)
	return err
}

func (repository *projectRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, project domain.Project) error {
	dto := toDTO(project)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO projects
				(id, slug, summary, description, creator_id, created_at, updated_at, deleted_at, started_at, finished_at, due_at, archived_at, priority_id, status_id, type_id)
			VALUES
				(?, ?, ?, ?, ?, ?, NULL, NULL, ?, ?, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		dto.Slug,
		dto.Summary,
		dto.Description,
		dto.CreatorID,
		dto.CreatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.ArchivedAt,
		dto.PriorityID,
		dto.StatusID,
		dto.TypeID,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *projectRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, project domain.Project) error {
	dto := toDTO(project)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE projects SET
				slug = ?,
				summary = ?,
				description = ?,
				updated_at = ?,
				started_at = ?,
				finished_at = ?,
				due_at = ?,
				archived_at = ?,
				priority_id = ?,
				status_id = ?,
				type_id = ?
			WHERE
				id = ?
			AND
				deleted_at IS NULL
        `,
		dto.Slug,
		dto.Summary,
		dto.Description,
		dto.UpdatedAt,
		dto.StartedAt,
		dto.FinishedAt,
		dto.DueAt,
		dto.ArchivedAt,
		dto.PriorityID,
		dto.StatusID,
		dto.TypeID,
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

func (repository *projectRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, deletedAt int64) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE projects SET
				deleted_at = ?
			WHERE id = ?
        `,
		deletedAt,
		projectID,
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

func (repository *projectRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) (domain.Project, error) {
	var dto projectDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                P.id,
				P.slug,
				P.summary,
				P.description,
				P.created_at,
				P.updated_at,
				P.deleted_at,
				P.started_at,
				P.finished_at,
				P.due_at,
				P.archived_at,
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
				IFNULL(PP.pages_count, 0) AS pages_count,
				IFNULL(PA.attachments_count, 0) AS attachments_count,
				IFNULL(PHO.history_operations_count, 0) AS history_operations_count,
				IFNULL(PT.tasks_count, 0) AS tasks_count
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
    			SELECT project_id, COUNT(*) AS pages_count
    			FROM project_pages
    			GROUP BY project_id
			) PP ON PP.project_id = P.id
			LEFT JOIN (
    			SELECT project_id, COUNT(*) AS attachments_count
    			FROM project_attachments
    			GROUP BY project_id
			) PA ON PA.project_id = P.id
			LEFT JOIN (
				SELECT project_id, COUNT(*) as history_operations_count
				FROM history_operations
				GROUP BY project_id
			) PHO ON PHO.project_id = P.id
			LEFT JOIN (
				SELECT project_id, COUNT(*) as tasks_count
				FROM tasks
				GROUP BY project_id
			) PT ON PHO.project_id = P.id
            WHERE
				P.id = ?
			GROUP
				BY P.id
        `,
		projectID,
	).Scan(
		&dto.ID,
		&dto.Slug,
		&dto.Summary,
		&dto.Description,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&dto.DeletedAt,
		&dto.StartedAt,
		&dto.FinishedAt,
		&dto.DueAt,
		&dto.ArchivedAt,
		&dto.StatusID,
		&dto.StatusName,
		&dto.StatusHexColor,
		&dto.PriorityID,
		&dto.PriorityName,
		&dto.PriorityHexColor,
		&dto.TypeID,
		&dto.TypeName,
		&dto.TypeHexColor,
		&dto.CreatorID,
		&dto.CreatorName,
		&dto.PermissionsCount,
		&dto.NotesCount,
		&dto.PagesCount,
		&dto.AttachmentsCount,
		&dto.HistoryOperationsCount,
		&dto.TasksCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Project{}, domain.NotFoundError
		}
		return domain.Project{}, err
	}
	return toDomain(dto), err

}

func (repository *projectRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.PagerQuery, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.PagerResult, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
                P.id,
				P.slug,
				P.summary,
				P.description,
				P.created_at,
				P.updated_at,
				P.deleted_at,
				P.started_at,
				P.finished_at,
				P.due_at,
				P.archived_at,
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
	sqlQueryVisibilityInnerJoins := ""
	var field string
	switch order.Field {
	case "slug":
		field = "P.slug COLLATE NOCASE"
	case "type":
		field = "PT.name COLLATE NOCASE"
	case "priority":
		field = "PP.name COLLATE NOCASE"
	case "status":
		field = "PS.item_index"
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
	case "archivedAt":
		field = "P.archived_at"
	case "createdBy":
		field = "U.name COLLATE NOCASE"
	default:
		field = "P.slug COLLATE NOCASE"
	}
	var sort string
	switch order.Direction {
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
	if filterDTO.ViewByUserID != nil && len(*filterDTO.ViewByUserID) > 0 {
		sqlQueryVisibilityInnerJoins = `
				INNER JOIN project_user_role ON project_user_role.project_id = P.id AND project_user_role.user_id = ?
				INNER JOIN roles ON roles.id = project_user_role.role_id AND roles.permissions_bitmask & ? = ?
			`
		filterArgs = append(filterArgs, *filterDTO.ViewByUserID)
		filterArgs = append(filterArgs, domain.PermissionViewProject)
		filterArgs = append(filterArgs, domain.PermissionViewProject)
	}
	if filterDTO.Slug != nil && len(*filterDTO.Slug) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.slug LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Slug+"%")
	}
	if filterDTO.Summary != nil && len(*filterDTO.Summary) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.summary LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Summary+"%")
	}
	if filterDTO.TypeID != nil && len(*filterDTO.TypeID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.type_id = ?")
		filterArgs = append(filterArgs, *filterDTO.TypeID)
	}
	if filterDTO.PriorityID != nil && len(*filterDTO.PriorityID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.priority_id = ?")
		filterArgs = append(filterArgs, *filterDTO.PriorityID)
	}
	if filterDTO.StatusID != nil && len(*filterDTO.StatusID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.status_id = ?")
		filterArgs = append(filterArgs, *filterDTO.StatusID)
	}
	if filterDTO.CreatedByUserID != nil && len(*filterDTO.CreatedByUserID) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "P.creator_id = ?")
		filterArgs = append(filterArgs, *filterDTO.CreatedByUserID)
	}
	if filterDTO.CreatedAt != nil {
		if filterDTO.CreatedAt.Empty != nil && *filterDTO.CreatedAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.created_at IS NULL")
		} else if filterDTO.CreatedAt.Filled != nil && *filterDTO.CreatedAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.created_at IS NOT NULL")
		} else {
			if filterDTO.CreatedAt.From != nil && *filterDTO.CreatedAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.created_at >= ?")
				filterArgs = append(filterArgs, filterDTO.CreatedAt.From)
			}
			if filterDTO.CreatedAt.To != nil && *filterDTO.CreatedAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.created_at <= ?")
				filterArgs = append(filterArgs, filterDTO.CreatedAt.To)
			}
		}
	}
	if filterDTO.UpdatedAt != nil {
		if filterDTO.UpdatedAt.Empty != nil && *filterDTO.UpdatedAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.updated_at IS NULL")
		} else if filterDTO.UpdatedAt.Filled != nil && *filterDTO.UpdatedAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.updated_at IS NOT NULL")
		} else {
			if filterDTO.UpdatedAt.From != nil && *filterDTO.UpdatedAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.updated_at >= ?")
				filterArgs = append(filterArgs, filterDTO.UpdatedAt.From)
			}
			if filterDTO.UpdatedAt.To != nil && *filterDTO.UpdatedAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.updated_at <= ?")
				filterArgs = append(filterArgs, filterDTO.UpdatedAt.To)
			}
		}
	}
	if filterDTO.StartedAt != nil {
		if filterDTO.StartedAt.Empty != nil && *filterDTO.StartedAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.started_at IS NULL")
		} else if filterDTO.StartedAt.Filled != nil && *filterDTO.StartedAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.started_at IS NOT NULL")
		} else {
			if filterDTO.StartedAt.From != nil && *filterDTO.StartedAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.started_at >= ?")
				filterArgs = append(filterArgs, filterDTO.StartedAt.From)
			}
			if filterDTO.StartedAt.To != nil && *filterDTO.StartedAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.started_at <= ?")
				filterArgs = append(filterArgs, filterDTO.StartedAt.To)
			}
		}
	}
	if filterDTO.FinishedAt != nil {
		if filterDTO.FinishedAt.Empty != nil && *filterDTO.FinishedAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.finished_at IS NULL")
		} else if filterDTO.FinishedAt.Filled != nil && *filterDTO.FinishedAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.finished_at IS NOT NULL")
		} else {
			if filterDTO.FinishedAt.From != nil && *filterDTO.FinishedAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.finished_at >= ?")
				filterArgs = append(filterArgs, filterDTO.FinishedAt.From)
			}
			if filterDTO.FinishedAt.To != nil && *filterDTO.FinishedAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.finished_at <= ?")
				filterArgs = append(filterArgs, filterDTO.FinishedAt.To)
			}
		}
	}
	if filterDTO.DueAt != nil {
		if filterDTO.DueAt.Empty != nil && *filterDTO.DueAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.due_at IS NULL")
		} else if filterDTO.DueAt.Filled != nil && *filterDTO.DueAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.due_at IS NOT NULL")
		} else {
			if filterDTO.DueAt.From != nil && *filterDTO.DueAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.due_at >= ?")
				filterArgs = append(filterArgs, filterDTO.DueAt.From)
			}
			if filterDTO.DueAt.To != nil && *filterDTO.DueAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.due_at <= ?")
				filterArgs = append(filterArgs, filterDTO.DueAt.To)
			}
		}
	}
	if filterDTO.ArchivedAt != nil {
		if filterDTO.ArchivedAt.Empty != nil && *filterDTO.ArchivedAt.Empty {
			sqlWhereConditions = append(sqlWhereConditions, "P.archived_at IS NULL")
		} else if filterDTO.ArchivedAt.Filled != nil && *filterDTO.ArchivedAt.Filled {
			sqlWhereConditions = append(sqlWhereConditions, "P.archived_at IS NOT NULL")
		} else {
			if filterDTO.ArchivedAt.From != nil && *filterDTO.ArchivedAt.From > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.archived_at >= ?")
				filterArgs = append(filterArgs, filterDTO.ArchivedAt.From)
			}
			if filterDTO.ArchivedAt.To != nil && *filterDTO.ArchivedAt.To > 0 {
				sqlWhereConditions = append(sqlWhereConditions, "P.archived_at <= ?")
				filterArgs = append(filterArgs, filterDTO.ArchivedAt.To)
			}
		}
	}
	/*
		if filterDTO.DeletedAt != nil {
			if filterDTO.DeletedAt.Empty != nil && *filterDTO.DeletedAt.Empty {
				sqlWhereConditions = append(sqlWhereConditions, "P.deleted_at IS NULL")
			} else if filterDTO.DeletedAt.Filled != nil && *filterDTO.DeletedAt.Filled {
				sqlWhereConditions = append(sqlWhereConditions, "P.deleted_at IS NOT NULL")
			} else {
				if filterDTO.DeletedAt.From != nil && *filterDTO.DeletedAt.From > 0 {
					sqlWhereConditions = append(sqlWhereConditions, "P.deleted_at >= ?")
					filterArgs = append(filterArgs, filterDTO.DeletedAt.From)
				}
				if filterDTO.DeletedAt.To != nil && *filterDTO.DeletedAt.To > 0 {
					sqlWhereConditions = append(sqlWhereConditions, "P.deleted_at <= ?")
					filterArgs = append(filterArgs, filterDTO.DeletedAt.To)
				}
			}
		}
	*/
	// TODO: FOR NOW always hide deleted projects, in future admins can list/restore deleted projects ?
	sqlWhereConditions = append(sqlWhereConditions, "P.deleted_at IS NULL")
	if len(sqlWhereConditions) > 0 {
		sqlWhere = " WHERE " + strings.Join(sqlWhereConditions, " AND ")
	}
	queryArgs = append(queryArgs, filterArgs...)
	var sqlLimit string
	if pager.Enabled {
		sqlLimit = " LIMIT ? OFFSET ? "
		queryArgs = append(queryArgs, pager.Limit(), pager.Offset())
	} else {
		sqlLimit = ""
	}
	sqlQuery = fmt.Sprintf("%s %s %s %s %s %s ", sqlQuery, sqlQueryInnerJoins, sqlQueryVisibilityInnerJoins, sqlWhere, sqlOrder, sqlLimit)
	//fmt.Println(sqlQuery)
	rows, err := dbExecutor.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.PagerResult{}, err
	}
	defer rows.Close()
	dtos := make([]projectDTO, 0)
	for rows.Next() {
		var dto projectDTO
		if err := rows.Scan(
			&dto.ID,
			&dto.Slug,
			&dto.Summary,
			&dto.Description,
			&dto.CreatedAt,
			&dto.UpdatedAt,
			&dto.DeletedAt,
			&dto.StartedAt,
			&dto.FinishedAt,
			&dto.DueAt,
			&dto.ArchivedAt,
			&dto.StatusID,
			&dto.StatusName,
			&dto.StatusHexColor,
			&dto.PriorityID,
			&dto.PriorityName,
			&dto.PriorityHexColor,
			&dto.TypeID,
			&dto.TypeName,
			&dto.TypeHexColor,
			&dto.CreatorID,
			&dto.CreatorName,
		); err != nil {
			return nil, browser.PagerResult{}, err
		}
		dtos = append(dtos, dto)
	}
	var totalResults int

	if pager.Enabled {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS total_projects
			FROM projects P
		`
		sqlCountQuery = fmt.Sprintf("%s %s %s", sqlCountQuery, sqlQueryVisibilityInnerJoins, sqlWhere)
		err = dbExecutor.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.PagerResult{}, err
		}
	} else {
		totalResults = len(dtos)
	}

	return toDomainArray(dtos), browser.NewPagerResult(pager, totalResults), nil
}

func (repository *projectRepository) GetCurrentProjects(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string, count uint) ([]domain.Project, error) {
	rows, err := dbExecutor.QueryContext(ctx,
		`
			SELECT
                P.id,
				P.slug,
				P.summary,
				P.description,
				P.created_at,
				P.updated_at,
				P.deleted_at,
				P.started_at,
				P.finished_at,
				P.due_at,
				P.archived_at,
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
			INNER JOIN project_user_role PUR ON PUR.project_id = P.id AND PUR.user_id = ?
			INNER JOIN project_priorities PP ON PP.id = P.priority_id
			INNER JOIN project_statuses PS ON PS.id = P.status_id
			INNER JOIN project_types PT ON PT.id = P.type_id
			INNER JOIN users U ON U.ID = P.creator_id
			WHERE P.archived_at IS NULL
			ORDER BY P.updated_at DESC, P.created_at DESC
			LIMIT ?
		`,
		userID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]projectDTO, 0)
	for rows.Next() {
		var dto projectDTO
		if err := rows.Scan(
			&dto.ID,
			&dto.Slug,
			&dto.Summary,
			&dto.Description,
			&dto.CreatedAt,
			&dto.UpdatedAt,
			&dto.DeletedAt,
			&dto.StartedAt,
			&dto.FinishedAt,
			&dto.DueAt,
			&dto.ArchivedAt,
			&dto.StatusID,
			&dto.StatusName,
			&dto.StatusHexColor,
			&dto.PriorityID,
			&dto.PriorityName,
			&dto.PriorityHexColor,
			&dto.TypeID,
			&dto.TypeName,
			&dto.TypeHexColor,
			&dto.CreatorID,
			&dto.CreatorName,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return toDomainArray(dtos), nil
}
