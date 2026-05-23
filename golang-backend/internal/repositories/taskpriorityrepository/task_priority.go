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
	Add(ctx context.Context, taskPriority taskPriorityDTO) error
	Update(ctx context.Context, taskPriority taskPriorityDTO) error
	Get(ctx context.Context, id string) (taskPriorityDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]taskPriorityDTO, browser.Result, error)
}

type taskPriorityRepository struct {
	database database.Database
}

func NewTaskPriorityRepository(database database.Database) TaskPriorityRepository {
	return &taskPriorityRepository{database: database}
}

func (taskPriorityRepository *taskPriorityRepository) Add(ctx context.Context, taskPriority taskPriorityDTO) error {
	_, err := taskPriorityRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO task_priorities (id, name, item_hex_color)
			VALUES (?, ?, ?)
        `,
		taskPriority.ID,
		taskPriority.Name,
		taskPriority.HexColor,
	)
	if err != nil {
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
			}
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			}
		}
	}
	return err
}

func (taskPriorityRepository *taskPriorityRepository) Update(ctx context.Context, projectPriority taskPriorityDTO) error {
	_, err := taskPriorityRepository.database.ExecContext(
		ctx,
		`
            UPDATE task_priorities SET
				name = ?,
				item_hex_color = ?
			WHERE id = ?
        `,
		projectPriority.Name,
		projectPriority.HexColor,
		projectPriority.ID,
	)
	if err != nil {
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
			}
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			}
		}
	}
	return err
}

func (taskPriorityRepository *taskPriorityRepository) Delete(ctx context.Context, id string) error {
	_, err := taskPriorityRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM task_priorities
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (taskPriorityRepository *taskPriorityRepository) Get(ctx context.Context, id string) (taskPriorityDTO, error) {
	var projectPriority taskPriorityDTO
	err := taskPriorityRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                TP.id, TP.name, TP.item_hex_color
            FROM task_priorities TP
            WHERE TP.id = ?
        `,
		id).Scan(&projectPriority.ID, &projectPriority.Name, &projectPriority.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return taskPriorityDTO{}, domain.NotFoundError
		}
		return taskPriorityDTO{}, err
	}
	return projectPriority, err
}

func (taskPriorityRepository *taskPriorityRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]taskPriorityDTO, browser.Result, error) {
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				TP.id, TP.name, TP.item_hex_color
			FROM task_priorities TP
    `
	var field string
	switch order.Field {
	case "name":
		field = "TP.name COLLATE NOCASE"
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
	if filter.Name != nil {
		sqlWhereConditions = append(sqlWhereConditions, "TP.name LIKE ?")
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
	rows, err := taskPriorityRepository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	taskPriorities := make([]taskPriorityDTO, 0)
	for rows.Next() {
		var taskPriority taskPriorityDTO
		if err := rows.Scan(
			&taskPriority.ID, &taskPriority.Name, &taskPriority.HexColor,
		); err != nil {
			return nil, browser.Result{}, err
		}
		taskPriorities = append(taskPriorities, taskPriority)
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
		err = taskPriorityRepository.database.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(taskPriorities)
	}

	return taskPriorities, browser.NewResult(pager, totalResults), nil
}
