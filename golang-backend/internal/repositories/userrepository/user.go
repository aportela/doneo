package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type UserRepository interface {
	Add(ctx context.Context, user domain.User, password string) error
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string, deletedAt int64) error
	UnDelete(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error)
}

type userRepository struct {
	database database.Database
}

func NewRepository(database database.Database) UserRepository {
	return &userRepository{database: database}
}

func (repository *userRepository) Add(ctx context.Context, user domain.User, password string) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dto := toDTO(user)
	dto.PasswordHash = string(hashedPasswordBytes)

	_, err = repository.database.ExecContext(
		ctx,
		`
            INSERT INTO users (id, email, name, password_hash, created_at, updated_at, deleted_at, permissions_bitmask)
			VALUES (?, ?, ?, ?, ?, NULL, NULL, ?)
        `,
		dto.ID,
		dto.Email,
		dto.Name,
		dto.PasswordHash,
		dto.CreatedAt,
		dto.PermissionsBitmask,
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

func (repository *userRepository) Update(ctx context.Context, user domain.User) error {
	dto := toDTO(user)
	if len(user.Password) > 0 {
		hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		dto.PasswordHash = string(hashedPasswordBytes)
	}
	var query string
	var args []any
	if dto.PasswordHash != "" {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				password_hash = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE id = ?`
		args = append(args, dto.Email, dto.Name, dto.PasswordHash, dto.UpdatedAt, dto.PermissionsBitmask, dto.ID)
	} else {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE id = ?`
		args = append(args, dto.Email, dto.Name, dto.UpdatedAt, dto.PermissionsBitmask, dto.ID)
	}
	_, err := repository.database.ExecContext(ctx, query, args...)
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

func (repository *userRepository) Delete(ctx context.Context, id string, deletedAt int64) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            UPDATE users SET
				deleted_at = ?
			WHERE id = ?
        `,
		deletedAt,
		id,
	)
	return err
}

func (repository *userRepository) UnDelete(ctx context.Context, id string) error {
	_, err := repository.database.ExecContext(
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

func (repository *userRepository) Purge(ctx context.Context, id string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM users
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (repository *userRepository) Get(ctx context.Context, id string) (domain.User, error) {
	var dto userDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
            FROM users U
            WHERE U.id = ?
        `,
		id).Scan(&dto.ID, &dto.Email, &dto.Name, &dto.PasswordHash, &dto.CreatedAt, &dto.UpdatedAt, &dto.DeletedAt, &dto.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}
	return toDomain(dto), err
}

func (repository *userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var dto userDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
            FROM users U
            WHERE U.email = ?
        `,
		email).Scan(&dto.ID, &dto.Email, &dto.Name, &dto.PasswordHash, &dto.CreatedAt, &dto.UpdatedAt, &dto.DeletedAt, &dto.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}
	return toDomain(dto), err
}

func (repository *userRepository) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	filterDTO := toFilterDTO(filter)
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
	if filterDTO.Name != nil && len(*filterDTO.Name) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "U.name LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Name+"%")
	}
	if filterDTO.Email != nil && len(*filterDTO.Email) > 0 {
		sqlWhereConditions = append(sqlWhereConditions, "U.email LIKE ?")
		filterArgs = append(filterArgs, "%"+*filterDTO.Email+"%")
	}
	if filterDTO.RequiredPermissionsBitmask != nil {
		sqlWhereConditions = append(sqlWhereConditions, "(U.permissions_bitmask & ?) = ?")
		filterArgs = append(filterArgs, filterDTO.RequiredPermissionsBitmask, filterDTO.RequiredPermissionsBitmask)
	}
	if filterDTO.ForbiddenPermissionsBitmask != nil {
		sqlWhereConditions = append(sqlWhereConditions, "(U.permissions_bitmask & ?) = 0")
		filterArgs = append(filterArgs, filterDTO.ForbiddenPermissionsBitmask)
	}
	if filterDTO.CreatedAt != nil {
		if filterDTO.CreatedAt.From != nil && *filterDTO.CreatedAt.From > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.created_at >= ?")
			filterArgs = append(filterArgs, filterDTO.CreatedAt.From)
		}
		if filterDTO.CreatedAt.To != nil && *filterDTO.CreatedAt.To > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.created_at <= ?")
			filterArgs = append(filterArgs, filterDTO.CreatedAt.To)
		}
	}
	if filterDTO.UpdatedAt != nil {
		if filterDTO.UpdatedAt.From != nil && *filterDTO.UpdatedAt.From > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.updated_at >= ?")
			filterArgs = append(filterArgs, filterDTO.UpdatedAt.From)
		}
		if filterDTO.UpdatedAt.To != nil && *filterDTO.UpdatedAt.To > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.updated_at <= ?")
			filterArgs = append(filterArgs, filterDTO.UpdatedAt.To)
		}
	}
	if filterDTO.DeletedAt != nil {
		if filterDTO.DeletedAt.From != nil && *filterDTO.DeletedAt.From > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.deleted_at >= ?")
			filterArgs = append(filterArgs, filterDTO.DeletedAt.From)
		}
		if filterDTO.DeletedAt.To != nil && *filterDTO.DeletedAt.To > 0 {
			sqlWhereConditions = append(sqlWhereConditions, "U.deleted_at <= ?")
			filterArgs = append(filterArgs, filterDTO.DeletedAt.To)
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
	rows, err := repository.database.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	users := make([]userDTO, 0)
	for rows.Next() {
		var user userDTO
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
		err = repository.database.QueryRowContext(
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

	return toDomainArray(users), browser.NewResult(pager, totalResults), nil
}
