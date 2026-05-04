package userrepository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type UserRepository interface {
	Add(ctx context.Context, user userDTO) error
	Update(ctx context.Context, user userDTO) error
	Delete(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (userDTO, error)
	GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (userDTO, error)
	Search(ctx context.Context) ([]userDTO, error)
}

type userRepository struct {
	database database.Database
}

func NewUserRepository(database database.Database) UserRepository {
	return &userRepository{database: database}
}

func (userRepository *userRepository) Add(ctx context.Context, user userDTO) error {

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

func (userRepository *userRepository) Update(ctx context.Context, user userDTO) error {
	var query string
	var args []interface{}
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

func (userRepository *userRepository) Get(ctx context.Context, id string) (userDTO, error) {
	var user userDTO
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
			return user, domain.ErrNotFound
		}
		return user, err
	}
	user.IsSuperUser = isSuperUser == 1
	return user, err
}

func (userRepository *userRepository) GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (userDTO, error) {
	var user userDTO
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
			return user, domain.ErrNotFound
		}
		return user, err
	}
	user.IsSuperUser = isSuperUser == 1
	return user, err
}

func (userRepository *userRepository) Search(ctx context.Context) ([]userDTO, error) {
	rows, err := userRepository.database.QueryContext(
		ctx,
		`
			SELECT
				U.id, U.email, U.name, U.created_at, U.updated_at, U.deleted_at, U.is_super_user
			FROM users U
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []userDTO
	for rows.Next() {
		var user userDTO
		var isSuperUser byte
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &isSuperUser); err != nil {
			return nil, err
		}
		user.IsSuperUser = isSuperUser == 1
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
