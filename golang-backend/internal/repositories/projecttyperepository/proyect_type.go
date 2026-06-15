package projecttyperepository

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

type ProjectTypeRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectType domain.ProjectType) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectType domain.ProjectType) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectTypeID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectTypeID string) (domain.ProjectType, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error)
}

type projectTypeRepository struct{}

func NewRepository() ProjectTypeRepository {
	return &projectTypeRepository{}
}

func (repository *projectTypeRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectType domain.ProjectType) error {
	dto := toDTO(projectType)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_types
				(id, name, item_hex_color)
			VALUES
				(?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.HexColor,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *projectTypeRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectType domain.ProjectType) error {
	dto := toDTO(projectType)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE project_types SET
				name = ?,
				item_hex_color = ?
			WHERE
				id = ?
        `,
		dto.Name,
		dto.HexColor,
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

func (repository *projectTypeRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectTypeID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_types
			WHERE
				id = ?
        `,
		projectTypeID,
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

func (repository *projectTypeRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectTypeID string) (domain.ProjectType, error) {
	var dto projectTypeDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                PT.id, PT.name, PT.item_hex_color
            FROM project_types PT
            WHERE
				PT.id = ?
        `,
		projectTypeID).Scan(&dto.ID, &dto.Name, &dto.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectType{}, domain.NotFoundError
		}
		return domain.ProjectType{}, err
	}
	return toDomain(dto), err
}

func (repository *projectTypeRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PT.id, PT.name, PT.item_hex_color
			FROM project_types PT
    `
	var field string
	switch order.Field {
	case "name":
		field = "PT.name COLLATE NOCASE"
	default:
		field = "PT.name COLLATE NOCASE"
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
		sqlWhereConditions = append(sqlWhereConditions, "PT.name LIKE ?")
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
	dtos := make([]projectTypeDTO, 0)
	for rows.Next() {
		var dto projectTypeDTO
		if err := rows.Scan(
			&dto.ID, &dto.Name, &dto.HexColor,
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
				COUNT(*) AS total_project_types
			FROM project_types PT
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
