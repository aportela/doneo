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
)

type UserRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, user domain.User, password string) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, user domain.User) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string, deletedAt int64) error
	UnDelete(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error
	Purge(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) (domain.User, error)
	GetByEmail(ctx context.Context, dbExecutor database.DatabaseExecutor, email string) (domain.User, error)
	SearchBase(ctx context.Context, dbExecutor database.DatabaseExecutor) ([]domain.UserBase, error)
	Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error)
}

type userRepository struct {
}

func NewRepository() UserRepository {
	return &userRepository{}
}

func (repository *userRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, user domain.User, password string) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	dto := toDTO(user)
	dto.PasswordHash = string(hashedPasswordBytes)

	_, err = dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO users
				(id, email, name, password_hash, created_at, updated_at, deleted_at, permissions_bitmask)
			VALUES
				(?, ?, ?, ?, ?, NULL, NULL, ?)
        `,
		dto.ID,
		dto.Email,
		dto.Name,
		dto.PasswordHash,
		dto.CreatedAt,
		dto.PermissionsBitmask,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *userRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, user domain.User) error {
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
			UPDATE users
			SET
				email = ?,
				name = ?,
				password_hash = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE
				id = ?`
		args = append(args, dto.Email, dto.Name, dto.PasswordHash, dto.UpdatedAt, dto.PermissionsBitmask, dto.ID)
	} else {
		query = `
			UPDATE users
			SET
				email = ?,
				name = ?,
				updated_at = ?,
				permissions_bitmask = ?
			WHERE
				id = ?`
		args = append(args, dto.Email, dto.Name, dto.UpdatedAt, dto.PermissionsBitmask, dto.ID)
	}
	result, err := dbExecutor.ExecContext(ctx, query, args...)
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

func (repository *userRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string, deletedAt int64) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE users
			SET
				deleted_at = ?
			WHERE
				id = ?
        `,
		deletedAt,
		userID,
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

func (repository *userRepository) UnDelete(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE users
			SET
				deleted_at = NULL
			WHERE
				id = ?
        `,
		userID,
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

func (repository *userRepository) Purge(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM users
			WHERE id = ?
        `,
		userID,
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

func (repository *userRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) (domain.User, error) {
	var dto userDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
            FROM users U
            WHERE
				U.id = ?
        `,
		userID).Scan(&dto.ID, &dto.Email, &dto.Name, &dto.PasswordHash, &dto.CreatedAt, &dto.UpdatedAt, &dto.DeletedAt, &dto.PermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, domain.NotFoundError
		}
		return domain.User{}, err
	}
	return toDomain(dto), err
}

func (repository *userRepository) GetByEmail(ctx context.Context, dbExecutor database.DatabaseExecutor, email string) (domain.User, error) {
	var dto userDTO
	err := dbExecutor.QueryRowContext(
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

func (repository *userRepository) SearchBase(ctx context.Context, dbExecutor database.DatabaseExecutor) ([]domain.UserBase, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
			U.id, U.email, U.name, U.created_at, U.updated_at, U.deleted_at, U.permissions_bitmask
			FROM users U
			ORDER BY U.name COLLATE NOCASE
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]userBaseDTO, 0)
	for rows.Next() {
		var dto userBaseDTO
		if err := rows.Scan(
			&dto.ID, &dto.Name,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toBaseDomainArray(dtos), nil
}

func (repository *userRepository) Search(ctx context.Context, dbExecutor database.DatabaseExecutor, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
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
	rows, err := dbExecutor.QueryContext(ctx, sqlQuery, queryArgs...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	dtos := make([]userDTO, 0)
	for rows.Next() {
		var dto userDTO
		if err := rows.Scan(&dto.ID, &dto.Email, &dto.Name, &dto.CreatedAt, &dto.UpdatedAt, &dto.DeletedAt, &dto.PermissionsBitmask); err != nil {
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
				COUNT(*) AS total_users
			FROM users U
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
