package projectpriorityrepository

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

type ProjectPriorityRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriority domain.ProjectPriority) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriority domain.ProjectPriority) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriorityID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriorityID string) (domain.ProjectPriority, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error)
}

type projectPriorityRepository struct{}

func NewRepository() ProjectPriorityRepository {
	return &projectPriorityRepository{}
}

func (repository *projectPriorityRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriority domain.ProjectPriority) error {
	dto := toDTO(projectPriority)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_priorities
				(id, name, item_hex_color, item_index)
			VALUES
				(?, ?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.HexColor,
		dto.Index,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *projectPriorityRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriority domain.ProjectPriority) error {
	dto := toDTO(projectPriority)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE project_priorities SET
				name = ?,
				item_hex_color = ?,
				item_index = ?
			WHERE
				id = ?
        `,
		dto.Name,
		dto.HexColor,
		dto.Index,
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

func (repository *projectPriorityRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriorityID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_priorities
			WHERE id = ?
        `,
		projectPriorityID,
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

func (repository *projectPriorityRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, projectPriorityID string) (domain.ProjectPriority, error) {
	var dto projectPriorityDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                PP.id, PP.name, PP.item_hex_color, PP.item_index
            FROM project_priorities PP
            WHERE PP.id = ?
        `,
		projectPriorityID).Scan(&dto.ID, &dto.Name, &dto.HexColor, &dto.Index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectPriority{}, domain.NotFoundError
		}
		return domain.ProjectPriority{}, err
	}
	return toDomain(dto), err
}

func (repository *projectPriorityRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PP.id, PP.name, PP.item_hex_color, PP.item_index
			FROM project_priorities PP
    `
	var field string
	switch order.Field {
	case "name":
		field = "PP.name COLLATE NOCASE"
	case "index":
		field = "PP.item_index"
	default:
		field = "PP.name COLLATE NOCASE"
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
		sqlWhereConditions = append(sqlWhereConditions, "PP.name LIKE ?")
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
	dtos := make([]projectPriorityDTO, 0)
	for rows.Next() {
		var dto projectPriorityDTO
		if err := rows.Scan(
			&dto.ID, &dto.Name, &dto.HexColor, &dto.Index,
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
				COUNT(*) AS total_project_priorities
			FROM project_priorities PP
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
