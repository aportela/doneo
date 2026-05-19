package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository interface {
	Add(ctx context.Context, user UserDTO) error
	Update(ctx context.Context, user UserDTO) error
	Patch(ctx context.Context, user UserDTO) error
	Delete(ctx context.Context, id string) error
	UnDelete(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (UserDTO, error)
	GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (UserDTO, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]UserDTO, browser.Result, error)
}

type userRepository struct {
	database database.Database
}

func NewUserRepository(database database.Database) UserRepository {
	return &userRepository{database: database}
}

func (userRepository *userRepository) Add(ctx context.Context, user UserDTO) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO users (id, email, name, password_hash, created_at, updated_at, deleted_at, permissions_bitmask)
			VALUES (?, ?, ?, ?, ?, NULL, NULL, ?)
        `,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.CreatedAt,
		user.PermissionsBitmask,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "users.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "users.email") {
				return &domain.AlreadyExistsError{Field: "email"}
			} else if strings.Contains(sqlErr.Error(), "users.id") {
				return &domain.AlreadyExistsError{Field: "id"}
			}
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(email)") {
				return &domain.ValidationError{Field: "email"}
			} else if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			}
		}
	}
	return err
}

func (userRepository *userRepository) Update(ctx context.Context, user UserDTO) error {
	var query string
	var args []any
	if user.PasswordHash != "" {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				password_hash = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE id = ?`
		args = append(args, user.Email, user.Name, user.PasswordHash, user.UpdatedAt, user.PermissionsBitmask, user.ID)
	} else {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE id = ?`
		args = append(args, user.Email, user.Name, user.UpdatedAt, user.PermissionsBitmask, user.ID)
	}
	_, err := userRepository.database.ExecContext(ctx, query, args...)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			if strings.Contains(sqlErr.Error(), "users.name") {
				return &domain.AlreadyExistsError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "users.email") {
				return &domain.AlreadyExistsError{Field: "email"}
			}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(name)") {
				return &domain.ValidationError{Field: "name"}
			} else if strings.Contains(sqlErr.Error(), "length(email)") {
				return &domain.ValidationError{Field: "email"}
			}
		}
	}
	return err
}

func (userRepository *userRepository) Patch(ctx context.Context, user UserDTO) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            UPDATE users
				SET
					deleted_at = ?
			WHERE id = ?
        `,
		user.DeletedAt,
		user.ID,
	)
	return err
}

func (userRepository *userRepository) Delete(ctx context.Context, id string) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            UPDATE users SET
				deleted_at = ?
			WHERE id = ?
        `,
		time.Now().UnixMilli(),
		id,
	)
	return err
}

func (userRepository *userRepository) UnDelete(ctx context.Context, id string) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            UPDATE users SET
				deleted_at = NULL
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (userRepository *userRepository) Purge(ctx context.Context, id string) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM users
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (userRepository *userRepository) Get(ctx context.Context, id string) (UserDTO, error) {
	var user UserDTO
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
            FROM users U
            WHERE U.id = ?
        `,
		id).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDTO{}, domain.NotFoundError
		}
		return UserDTO{}, err
	}
	return user, err
}

func (userRepository *userRepository) GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (UserDTO, error) {
	var user UserDTO
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
            FROM users U
            WHERE U.email = ?
        `,
		email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDTO{}, domain.NotFoundError
		}
		return UserDTO{}, err
	}
	return user, err
}

func (userRepository *userRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter searchFilterDTO) ([]UserDTO, browser.Result, error) {
	var filterArgs []any
	var queryArgs []any
	sqlQuery := `
		SELECT
			U.id, U.email, U.name, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
		FROM users U
	`
	var field string
	switch order.Field {
	case "name":
		field = "U.name COLLATE NOCASE"
	case "email":
		field = "U.email COLLATE NOCASE"
	case "createdAt":
		field = "U.created_at"
	case "updatedAt":
		field = "U.updated_at"
	case "deletedAt":
		field = "U.deleted_at"
	case "isSuperUser":
		field = "U.permissions_bitmask"
	default:
		field = "U.name COLLATE NOCASE"
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
		sqlWhereConditions = append(sqlWhereConditions, "U.name LIKE ?")
		filterArgs = append(filterArgs, "%"+*filter.Name+"%")
	}
	if filter.Email != nil {
		sqlWhereConditions = append(sqlWhereConditions, "U.email LIKE ?")
		filterArgs = append(filterArgs, "%"+*filter.Email+"%")
	}
	if filter.RequiredPermissionsBitmask != nil {
		sqlWhereConditions = append(sqlWhereConditions, "(U.permissions_bitmask & ?) = ?")
		filterArgs = append(filterArgs, filter.RequiredPermissionsBitmask, filter.RequiredPermissionsBitmask)
	}
	if filter.ForbiddenPermissionsBitmask != nil {
		sqlWhereConditions = append(sqlWhereConditions, "(U.permissions_bitmask & ?) = 0")
		filterArgs = append(filterArgs, filter.ForbiddenPermissionsBitmask)
	}
	if filter.CreatedAt != nil {
		if filter.CreatedAt.From != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.created_at >= ?")
			filterArgs = append(filterArgs, filter.CreatedAt.From)
		}
		if filter.CreatedAt.To != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.created_at <= ?")
			filterArgs = append(filterArgs, filter.CreatedAt.To)
		}
	}
	if filter.UpdatedAt != nil {
		if filter.UpdatedAt.From != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.updated_at >= ?")
			filterArgs = append(filterArgs, filter.UpdatedAt.From)
		}
		if filter.UpdatedAt.To != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.updated_at <= ?")
			filterArgs = append(filterArgs, filter.UpdatedAt.To)
		}
	}
	if filter.DeletedAt != nil {
		if filter.DeletedAt.From != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.deleted_at >= ?")
			filterArgs = append(filterArgs, filter.DeletedAt.From)
		}
		if filter.DeletedAt.To != nil {
			sqlWhereConditions = append(sqlWhereConditions, "U.deleted_at <= ?")
			filterArgs = append(filterArgs, filter.DeletedAt.To)
		}
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
	rows, err := userRepository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	users := make([]UserDTO, 0)
	for rows.Next() {
		var user UserDTO
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.PermissionsBitmask); err != nil {
			return nil, browser.Result{}, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, browser.Result{}, err
	}

	var totalResults int

	if pager.Enabled() {
		sqlCountQuery := `
			SELECT
				COUNT(*) AS total_users
			FROM users U
		`
		sqlCountQuery = fmt.Sprintf("%s %s", sqlCountQuery, sqlWhere)
		err = userRepository.database.QueryRowContext(
			ctx,
			sqlCountQuery,
			filterArgs...,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(users)
	}

	return users, browser.NewResult(pager, totalResults), nil
}
