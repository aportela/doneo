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
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type TaskPriorityRepository interface {
	Add(ctx context.Context, taskPriority domain.TaskPriority) error
	Update(ctx context.Context, taskPriority domain.TaskPriority) error
	Get(ctx context.Context, id string) (domain.TaskPriority, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error)
}

type taskPriorityRepository struct {
	db database.Database
}

func NewRepository(db database.Database) TaskPriorityRepository {
	return &taskPriorityRepository{db: db}
}

func (repository *taskPriorityRepository) Add(ctx context.Context, taskPriority domain.TaskPriority) error {
	dto := toDTO(taskPriority)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO task_priorities (id, name, item_hex_color, item_index)
			VALUES (?, ?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.HexColor,
		dto.Index,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "task_priorities.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "task_priorities.id") {
				return &domain.AlreadyExistsError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "task_priorities.item_index") {
				return &domain.AlreadyExistsError{Field: "index"}
			}
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskPriorityRepository) Update(ctx context.Context, taskPriority domain.TaskPriority) error {
	dto := toDTO(taskPriority)
	_, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE task_priorities SET
				name = ?,
				item_hex_color = ?,
				item_index = ?
			WHERE id = ?
        `,
		dto.Name,
		dto.HexColor,
		dto.Index,
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
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "task_priorities.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "task_priorities.id") {
				return &domain.AlreadyExistsError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "task_priorities.item_index") {
				return &domain.AlreadyExistsError{Field: "index"}
			}
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskPriorityRepository) Delete(ctx context.Context, id string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM task_priorities
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (repository *taskPriorityRepository) Get(ctx context.Context, id string) (domain.TaskPriority, error) {
	var dto taskPriorityDTO
	err := repository.db.QueryRowContext(
		ctx,
		`
            SELECT
                TP.id, TP.name, TP.item_hex_color, TP.item_index
            FROM task_priorities TP
            WHERE TP.id = ?
        `,
		id).Scan(&dto.ID, &dto.Name, &dto.HexColor, &dto.Index)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskPriority{}, domain.NotFoundError
		}
		return domain.TaskPriority{}, err
	}
	return toDomain(dto), err
}

func (repository *taskPriorityRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchTaskPrioritiesFilter) ([]domain.TaskPriority, browser.Result, error) {
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
	rows, err := repository.db.QueryContext(ctx, sqlQuery, queryArgs...)
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
		err = repository.db.QueryRowContext(
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
