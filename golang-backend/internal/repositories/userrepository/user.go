package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
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
	Search(ctx context.Context, pager browser.Params) ([]UserDTO, browser.Result, error)
}

type userRepository struct {
	database database.Database
}

func NewUserRepository(database database.Database) UserRepository {
	return &userRepository{database: database}
}

func (userRepository *userRepository) Add(ctx context.Context, user UserDTO) error {
	adminFlag := 0
	if user.IsSuperUser {
		adminFlag = 1
	}
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO users (id, email, name, password_hash, created_at, updated_at, deleted_at, is_super_user)
			VALUES (?, ?, ?, ?, ?, NULL, NULL, ?)
        `,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.CreatedAt,
		adminFlag,
	)
	return err
}

func (userRepository *userRepository) Update(ctx context.Context, user UserDTO) error {
	var query string
	var args []any
	adminFlag := 0
	if user.IsSuperUser {
		adminFlag = 1
	}
	if user.PasswordHash != "" {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				password_hash = ?,
				updated_at = ?,
				is_super_user = ?
			WHERE id = ?`
		args = append(args, user.Email, user.Name, user.PasswordHash, user.UpdatedAt, adminFlag, user.ID)
	} else {
		query = `
			UPDATE users SET
				email = ?,
				name = ?,
				updated_at = ?,
				is_super_user = ?
			WHERE id = ?`
		args = append(args, user.Email, user.Name, user.UpdatedAt, adminFlag, user.ID)
	}
	_, err := userRepository.database.ExecContext(ctx, query, args...)
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
	var isSuperUser byte
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.is_super_user
            FROM users U
            WHERE U.id = ?
        `,
		id).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &isSuperUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDTO{}, domain.ErrNotFound
		}
		return UserDTO{}, err
	}
	user.IsSuperUser = isSuperUser == 1
	return user, err
}

func (userRepository *userRepository) GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (UserDTO, error) {
	var user UserDTO
	var isSuperUser byte
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.deleted_at, U.is_super_user
            FROM users U
            WHERE U.email = ?
        `,
		email).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &isSuperUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserDTO{}, domain.ErrNotFound
		}
		return UserDTO{}, err
	}
	user.IsSuperUser = isSuperUser == 1
	return user, err
}

func (userRepository *userRepository) Search(ctx context.Context, pager browser.Params) ([]UserDTO, browser.Result, error) {
	var args []any
	query := `
		SELECT
			U.id, U.email, U.name, U.created_at, U.updated_at, U.deleted_at, U.is_super_user
		FROM users U
	`
	query += " ORDER BY U.name COLLATE NOCASE "
	if pager.Enabled() {
		query += " LIMIT ? OFFSET ? "
		args = append(args, pager.Limit(), pager.Offset())
	}
	rows, err := userRepository.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, browser.Result{}, err
	}
	defer rows.Close()
	var users []UserDTO
	for rows.Next() {
		var user UserDTO
		var isSuperUser byte
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &isSuperUser); err != nil {
			return nil, browser.Result{}, err
		}
		user.IsSuperUser = isSuperUser == 1
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, browser.Result{}, err
	}

	var totalResults int

	if pager.Enabled() {
		err = userRepository.database.QueryRowContext(
			ctx,
			`
				SELECT
					COUNT(*) AS total_users
				FROM users
			`,
		).Scan(&totalResults)

		if err != nil {
			return nil, browser.Result{}, err
		}
	} else {
		totalResults = len(users)
	}

	return users, browser.NewResult(pager, totalResults), nil
}
