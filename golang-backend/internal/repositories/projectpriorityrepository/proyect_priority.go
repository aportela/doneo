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
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectPriorityRepository interface {
	Add(ctx context.Context, projectPriority projectPriorityDTO) error
	Update(ctx context.Context, projectPriority projectPriorityDTO) error
	Get(ctx context.Context, id string) (projectPriorityDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectPriorityDTO, browser.Result, error)
}

type projectPriorityRepository struct {
	database database.Database
}

func NewProjectPriorityRepository(database database.Database) ProjectPriorityRepository {
	return &projectPriorityRepository{database: database}
}

func (projectPriorityRepository *projectPriorityRepository) Add(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_priorities (id, name, item_hex_color)
			VALUES (?, ?, ?)
        `,
		projectPriority.ID,
		projectPriority.Name,
		projectPriority.HexColor,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "project_priorities.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "project_priorities.id") {
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

func (projectPriorityRepository *projectPriorityRepository) Update(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_priorities SET
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
			if strings.Contains(sqlErr.Error(), "project_priorities.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "project_priorities.id") {
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

func (projectPriorityRepository *projectPriorityRepository) Delete(ctx context.Context, id string) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_priorities
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (projectPriorityRepository *projectPriorityRepository) Get(ctx context.Context, id string) (projectPriorityDTO, error) {
	var projectPriority projectPriorityDTO
	err := projectPriorityRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PP.id, PP.name, PP.item_hex_color
            FROM project_priorities PP
            WHERE PP.id = ?
        `,
		id).Scan(&projectPriority.ID, &projectPriority.Name, &projectPriority.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectPriorityDTO{}, domain.NotFoundError
		}
		return projectPriorityDTO{}, err
	}
	return projectPriority, err
}

func (projectPriorityRepository *projectPriorityRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]projectPriorityDTO, browser.Result, error) {
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PP.id, PP.name, PP.item_hex_color
			FROM project_priorities PP
    `
	var field string
	switch order.Field {
	case "name":
		field = "PP.name COLLATE NOCASE"
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
	if filter.Name != nil {
		sqlWhereConditions = append(sqlWhereConditions, "PP.name LIKE ?")
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
	rows, err := projectPriorityRepository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	projectPriorities := make([]projectPriorityDTO, 0)
	for rows.Next() {
		var projectPriority projectPriorityDTO
		if err := rows.Scan(
			&projectPriority.ID, &projectPriority.Name, &projectPriority.HexColor,
		); err != nil {
			return nil, browser.Result{}, err
		}
		projectPriorities = append(projectPriorities, projectPriority)
	}
	if err := rows.Err(); err != nil {
		return nil, browser.Result{}, err
	}

	var totalResults int

	if pager.Enabled() {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS project_priorities
			FROM project_priorities PP
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = projectPriorityRepository.database.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(projectPriorities)
	}

	return projectPriorities, browser.NewResult(pager, totalResults), nil
}
