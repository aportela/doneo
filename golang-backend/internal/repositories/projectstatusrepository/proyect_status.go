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
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectStatusRepository interface {
	Add(ctx context.Context, projectStatus domain.ProjectStatus) error
	Update(ctx context.Context, projectStatus domain.ProjectStatus) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectStatus, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error)
}

type projectStatusRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectStatusRepository {
	return &projectStatusRepository{database: database}
}

func (repository *projectStatusRepository) Add(ctx context.Context, projectStatus domain.ProjectStatus) error {
	dto := toDTO(projectStatus)
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_statuses (id, name, item_hex_color)
			VALUES (?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.HexColor,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "project_statuses.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "project_statuses.id") {
				return &domain.AlreadyExistsError{Field: "id"}
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

func (repository *projectStatusRepository) Update(ctx context.Context, projectStatus domain.ProjectStatus) error {
	dto := toDTO(projectStatus)
	_, err := repository.database.ExecContext(
		ctx,
		`
            UPDATE project_statuses SET
				name = ?,
				item_hex_color = ?
			WHERE id = ?
        `,
		dto.Name,
		dto.HexColor,
		dto.ID,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "project_statuses.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "project_statuses.id") {
				return &domain.AlreadyExistsError{Field: "id"}
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

func (repository *projectStatusRepository) Delete(ctx context.Context, id string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_statuses
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (repository *projectStatusRepository) Get(ctx context.Context, id string) (domain.ProjectStatus, error) {
	var dto projectStatusDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PS.id, PS.name, PS.item_hex_color
            FROM project_statuses PS
            WHERE PS.id = ?
        `,
		id).Scan(&dto.ID, &dto.Name, &dto.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectStatus{}, domain.NotFoundError
		}
		return domain.ProjectStatus{}, err
	}
	return toDomain(dto), err
}

func (repository *projectStatusRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
			SELECT
				PS.id, PS.name, PS.item_hex_color
			FROM project_statuses PS
    `
	var field string
	switch order.Field {
	case "name":
		field = "PS.name COLLATE NOCASE"
	default:
		field = "PS.name COLLATE NOCASE"
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
	rows, err := repository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	dtos := make([]projectStatusDTO, 0)
	for rows.Next() {
		var dto projectStatusDTO
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
				COUNT(*) AS total_project_statuses
			FROM project_statuses PT
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = repository.database.QueryRowContext(
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
