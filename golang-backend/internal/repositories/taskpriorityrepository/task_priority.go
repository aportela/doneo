package taskpriorityrepository

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

type TaskPriorityRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriority domain.TaskPriority) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriority domain.TaskPriority) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriorityID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriorityID string) (domain.TaskPriority, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error)
}

type taskPriorityRepository struct{}

func NewRepository() TaskPriorityRepository {
	return &taskPriorityRepository{}
}

func (repository *taskPriorityRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriority domain.TaskPriority) error {
	dto := toDTO(taskPriority)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_priorities
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

func (repository *taskPriorityRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriority domain.TaskPriority) error {
	dto := toDTO(taskPriority)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE task_priorities SET
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

func (repository *taskPriorityRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriorityID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_priorities
			WHERE
				id = ?
        `,
		taskPriorityID,
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

func (repository *taskPriorityRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskPriorityID string) (domain.TaskPriority, error) {
	var dto taskPriorityDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                TP.id, TP.name, TP.item_hex_color, TP.item_index
            FROM task_priorities TP
            WHERE
				TP.id = ?
        `,
		taskPriorityID).Scan(&dto.ID, &dto.Name, &dto.HexColor, &dto.Index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskPriority{}, domain.NotFoundError
		}
		return domain.TaskPriority{}, err
	}
	return toDomain(dto), err
}

func (repository *taskPriorityRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				TP.id, TP.name, TP.item_hex_color, TP.item_index
			FROM task_priorities TP
    `
	var field string
	switch order.Field {
	case "name":
		field = "TP.name COLLATE NOCASE"
	case "index":
		field = "TP.item_index"
	default:
		field = "TP.name COLLATE NOCASE"
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
		sqlWhereConditions = append(sqlWhereConditions, "TP.name LIKE ?")
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
	dtos := make([]taskPriorityDTO, 0)
	for rows.Next() {
		var dto taskPriorityDTO
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
				COUNT(*) AS task_priorities
			FROM task_priorities TP
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
