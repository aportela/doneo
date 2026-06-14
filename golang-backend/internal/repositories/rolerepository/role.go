package rolerepository

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

type RoleRepository interface {
	Add(ctx context.Context, role domain.Role) error
	Update(ctx context.Context, role domain.Role) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Role, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error)
}

type roleRepository struct {
	db database.Database
}

func NewRepository(db database.Database) RoleRepository {
	return &roleRepository{db: db}
}

func (repository *roleRepository) Add(ctx context.Context, role domain.Role) error {
	dto := toDTO(role)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO roles (id, name, permissions_bitmask)
			VALUES (?, ?, ?)
        `,
		dto.ID,
		dto.Name,
		dto.PermissionsBitmask,
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
			if strings.Contains(sqlErr.Error(), "roles.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "roles.id") {
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

func (repository *roleRepository) Update(ctx context.Context, role domain.Role) error {
	dto := toDTO(role)
	_, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE roles SET
				name = ?,
				permissions_bitmask = ?
			WHERE id = ?
        `,
		dto.Name,
		dto.PermissionsBitmask,
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
			if strings.Contains(sqlErr.Error(), "roles.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			}
			return err
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *roleRepository) Delete(ctx context.Context, id string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM roles
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (repository *roleRepository) Get(ctx context.Context, id string) (domain.Role, error) {
	var dto roleDTO
	err := repository.db.QueryRowContext(
		ctx,
		`
            SELECT
                R.id, R.name, R.permissions_bitmask
            FROM roles R
            WHERE R.id = ?
        `,
		id).Scan(&dto.ID, &dto.Name, &dto.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Role{}, domain.NotFoundError
		}
		return domain.Role{}, err
	}
	return toDomain(dto), err
}

func (repository *roleRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
			R.id, R.name, R.permissions_bitmask
		FROM roles R
	`
	var field string
	switch order.Field {
	case "name":
		field = "R.name COLLATE NOCASE"
	default:
		field = "R.name COLLATE NOCASE"
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
		sqlWhereConditions = append(sqlWhereConditions, "R.name LIKE ?")
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
	dtos := make([]roleDTO, 0)
	for rows.Next() {
		var dto roleDTO
		if err := rows.Scan(&dto.ID, &dto.Name, &dto.PermissionsBitmask); err != nil {
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
				COUNT(*) AS total_roles
			FROM roles R
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
