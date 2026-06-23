package taskstatusrepository

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

type TaskStatusRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatus domain.TaskStatus) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatus domain.TaskStatus) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatusID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatusID string) (domain.TaskStatus, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error)
}

type taskStatusRepository struct{}

func NewRepository() TaskStatusRepository {
	return &taskStatusRepository{}
}

func (repository *taskStatusRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatus domain.TaskStatus) error {
	dto := toDTO(taskStatus)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_statuses
				(id, name, item_hex_color, item_index, flags_bitmask)
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

func (repository *taskStatusRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatus domain.TaskStatus) error {
	dto := toDTO(taskStatus)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE task_statuses SET
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

func (repository *taskStatusRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatusID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_statuses
			WHERE
				id = ?
        `,
		taskStatusID,
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

func (repository *taskStatusRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskStatusID string) (domain.TaskStatus, error) {
	var dto taskStatusDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                TS.id, TS.name, TS.item_hex_color, TS.item_index, TS.flags_bitmask
            FROM task_statuses TS
            WHERE
				TS.id = ?
        `,
		taskStatusID).Scan(&dto.ID, &dto.Name, &dto.HexColor, &dto.Index, &dto.FlagsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskStatus{}, domain.NotFoundError
		}
		return domain.TaskStatus{}, err
	}
	return toDomain(dto), err
}

func (repository *taskStatusRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchTaskStatusesFilter) ([]domain.TaskStatus, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				TS.id, TS.name, TS.item_hex_color, TS.item_index, TS.flags_bitmask
			FROM task_statuses TS
    `
	var field string
	switch order.Field {
	case "name":
		field = "TS.name COLLATE NOCASE"
	case "index":
		field = "TS.item_index"
	default:
		field = "TS.item_index"
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
		sqlWhereConditions = append(sqlWhereConditions, "TS.name LIKE ?")
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
	dtos := make([]taskStatusDTO, 0)
	for rows.Next() {
		var dto taskStatusDTO
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
				COUNT(*) AS total_task_statuses
			FROM task_statuses TS
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
