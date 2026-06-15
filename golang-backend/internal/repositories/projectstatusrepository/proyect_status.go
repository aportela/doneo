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
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatus domain.ProjectStatus) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatus domain.ProjectStatus) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatusID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatusID string) (domain.ProjectStatus, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error)
}

type projectStatusRepository struct {
	db database.Database
}

func NewRepository() ProjectStatusRepository {
	return &projectStatusRepository{}
}

func (repository *projectStatusRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatus domain.ProjectStatus) error {
	dto := toDTO(projectStatus)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO
				project_statuses (id, name, item_hex_color, item_index, flags_bitmask)
			VALUES
				(?, ?, ?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.HexColor,
		dto.Index,
		dto.FlagsBitmask,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *projectStatusRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatus domain.ProjectStatus) error {
	dto := toDTO(projectStatus)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE project_statuses SET
				name = ?,
				item_hex_color = ?,
				item_index = ?,
				flags_bitmask = ?
			WHERE
				id = ?
        `,
		dto.Name,
		dto.HexColor,
		dto.Index,
		dto.FlagsBitmask,
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

func (repository *projectStatusRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatusID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_statuses
			WHERE
				id = ?
        `,
		projectStatusID,
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

func (repository *projectStatusRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectStatusID string) (domain.ProjectStatus, error) {
	var dto projectStatusDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                PS.id, PS.name, PS.item_hex_color, PS.item_index, PS.flags_bitmask
            FROM project_statuses PS
            WHERE PS.id = ?
        `,
		projectStatusID).Scan(&dto.ID, &dto.Name, &dto.HexColor, &dto.Index, &dto.FlagsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectStatus{}, domain.NotFoundError
		}
		return domain.ProjectStatus{}, err
	}
	return toDomain(dto), err
}

func (repository *projectStatusRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PS.id, PS.name, PS.item_hex_color, PS.item_index, PS.flags_bitmask
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
	if filterDTO.Name != nil && len(*filterDTO.Name) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "PS.name LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Name+"%")
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
	rows, err := dbExecutor.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	dtos := make([]projectStatusDTO, 0)
	for rows.Next() {
		var dto projectStatusDTO
		if err := rows.Scan(
			&dto.ID, &dto.Name, &dto.HexColor, &dto.Index, &dto.FlagsBitmask,
		); err != nil {
			return nil, browser.Result{}, err
		}
		dtos = append(dtos, dto)
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
