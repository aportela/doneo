package projectstatusrepository

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

type ProjectStatusRepository interface {
	Add(ctx context.Context, projectStatus projectStatusDTO) error
	Update(ctx context.Context, projectStatus projectStatusDTO) error
	Get(ctx context.Context, id string) (projectStatusDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectStatusDTO, browser.Result, error)
}

type projectStatusRepository struct {
	database database.Database
}

func NewProjectStatusRepository(database database.Database) ProjectStatusRepository {
	return &projectStatusRepository{database: database}
}

func (projectStatusRepository *projectStatusRepository) Add(ctx context.Context, projectStatus projectStatusDTO) error {
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_statuses (id, name, item_hex_color, item_index)
			VALUES (?, ?, ?, COALESCE((SELECT MAX(item_index) FROM project_statuses), 0) + 1)
        `,
		projectStatus.ID,
		projectStatus.Name,
		projectStatus.HexColor,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Update(ctx context.Context, projectStatus projectStatusDTO) error {
	/*
		`
			WITH moved AS (
				SELECT id, item_index AS old_index
				FROM project_statuses
				WHERE id = ?
			)
			UPDATE project_statuses
			SET item_index = CASE
				WHEN id = (SELECT id FROM moved) THEN ?  -- nuevo índice
				WHEN (SELECT old_index FROM moved) > ?
					AND item_index >= ? AND item_index < (SELECT old_index FROM moved) THEN item_index + 1
				WHEN (SELECT old_index FROM moved) < ?
					AND item_index <= ? AND item_index > (SELECT old_index FROM moved) THEN item_index - 1
				ELSE item_index
			END
			WHERE item_index BETWEEN MIN(?, (SELECT old_index FROM moved)) AND MAX(?, (SELECT old_index FROM moved));
		`
	*/
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_statuses SET
				name = ?,
				item_hex_color = ?,
				item_index = ?
			WHERE id = ?
        `,
		projectStatus.Name,
		projectStatus.HexColor,
		projectStatus.Index,
		projectStatus.ID,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Delete(ctx context.Context, id string) error {
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_statuses
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Get(ctx context.Context, id string) (projectStatusDTO, error) {
	var projectStatus projectStatusDTO
	err := projectStatusRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PS.id, PS.name, PS.item_hex_color, PS.item_index
            FROM project_statuses PS
            WHERE PS.id = ?
        `,
		id).Scan(&projectStatus.ID, &projectStatus.Name, &projectStatus.HexColor, &projectStatus.Index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectStatusDTO{}, domain.ErrNotFound
		}
		return projectStatusDTO{}, err
	}
	return projectStatus, err
}

func (projectStatusRepository *projectStatusRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectStatusDTO, browser.Result, error) {
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PS.id, PS.name, PS.item_hex_color, PS.item_index
			FROM project_statuses PS
    `
	var field string
	switch order.Field {
	case "name":
		field = "PS.name COLLATE NOCASE"
	case "index":
		field = "PS.item_index"
	default:
		field = "PS.item_index"
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
	if filter.Name != nil {
		sqlWhereConditions = append(sqlWhereConditions, "PS.name LIKE ?")
		filterArgs = append(filterArgs, "%"+*filter.Name+"%")
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
	sqlQuery = fmt.Sprintf("%s %s %s %s ", sqlQuery, sqlWhere, sqlOrder, sqlLimit)
	rows, err := projectStatusRepository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	projectStatuses := make([]projectStatusDTO, 0)
	for rows.Next() {
		var projectStatus projectStatusDTO
		if err := rows.Scan(
			&projectStatus.ID, &projectStatus.Name, &projectStatus.HexColor, &projectStatus.Index,
		); err != nil {
			return nil, browser.Result{}, err
		}
		projectStatuses = append(projectStatuses, projectStatus)
	}
	if err := rows.Err(); err != nil {
		return nil, browser.Result{}, err
	}

	var totalResults int

	if pager.Enabled() {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS total_project_statuses
			FROM project_statuses PT
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = projectStatusRepository.database.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(projectStatuses)
	}

	return projectStatuses, browser.NewResult(pager, totalResults), nil
}
