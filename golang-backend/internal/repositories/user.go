package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

type UserRepository interface {
	Add(ctx context.Context, user domain.User) error
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	GetById(ctx context.Context, id string) (domain.User, error)
	GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (domain.User, error)
	Search(ctx context.Context) ([]domain.User, error)
}

type userRepository struct {
	database database.Database
}

func NewUserRepository(database database.Database) UserRepository {
	return &userRepository{database: database}
}

func getAvatar(userId string) string {
	return "https://i.pravatar.cc/48?id=" + userId
}

func (userRepository *userRepository) Add(ctx context.Context, user domain.User) error {

	adminFlag := 0
	if user.IsSuperUser {
		adminFlag = 1
	}
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO users (id, email, name, password_hash, created_at, updated_at, is_super_user)
			VALUES (?, ?, ?, ?, ?, NULL, ?)
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

func (userRepository *userRepository) Update(ctx context.Context, user domain.User) error {

	var query string
	var args []interface{}
	if user.PasswordHash != nil {
		query = `UPDATE users SET email = ?, name = ?, password_hash = ?, updated_at = ? WHERE id = ?`
		args = append(args, user.Email, user.Name, user.PasswordHash, user.UpdatedAt, user.ID)
	} else {
		query = `UPDATE users SET email = ?, name = ?, updated_at = ? WHERE id = ?`
		args = append(args, user.Email, user.Name, user.UpdatedAt, user.ID)
	}
	_, err := userRepository.database.ExecContext(ctx, query, args...)
	return err
}

func (userRepository *userRepository) Delete(ctx context.Context, id string) error {
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

func (userRepository *userRepository) GetById(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	var updatedAt sql.NullInt64
	var isSuperUser sql.NullByte
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.created_at, U.updated_at, U.is_super_user
            FROM users U
            WHERE U.id = ?
        `,
		id).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &updatedAt, &isSuperUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrNotFound
		}
		return user, err
	}
	user.UpdatedAt = utils.SQLInt64Ptr(updatedAt)
	user.IsSuperUser = isSuperUser.Valid && isSuperUser.Byte == 1
	user.AvatarURL = getAvatar(user.ID)
	return user, err
}

func (userRepository *userRepository) GetByEmailForVerifyCredentials(ctx context.Context, email string, password string) (domain.User, error) {
	var user domain.User
	var updatedAt sql.NullInt64
	var isSuperUser sql.NullByte
	var passwordHash string
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.password_hash, U.created_at, U.updated_at, U.is_super_user
            FROM users U
            WHERE U.email = ?
        `,
		email).Scan(&user.ID, &user.Email, &user.Name, &passwordHash, &user.CreatedAt, &updatedAt, &isSuperUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrNotFound
		}
		return user, err
	}
	user.PasswordHash = &passwordHash
	user.UpdatedAt = utils.SQLInt64Ptr(updatedAt)
	user.IsSuperUser = isSuperUser.Valid && isSuperUser.Byte == 1
	user.AvatarURL = getAvatar(user.ID)
	return user, err
}

func (userRepository *userRepository) Search(ctx context.Context) ([]domain.User, error) {
	rows, err := userRepository.database.QueryContext(
		ctx,
		`
			SELECT
				U.id, U.email, U.name, U.created_at, U.updated_at, U.is_super_user
			FROM users U
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []domain.User
	for rows.Next() {
		var user domain.User
		var updatedAt sql.NullInt64
		var isSuperUser sql.NullByte
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &updatedAt, &isSuperUser); err != nil {
			return nil, err
		}

		user.UpdatedAt = utils.SQLInt64Ptr(updatedAt)
		user.IsSuperUser = isSuperUser.Valid && isSuperUser.Byte == 1
		user.AvatarURL = getAvatar(user.ID)
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
