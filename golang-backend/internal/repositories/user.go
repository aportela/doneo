package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/models"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	database database.Database
}

func NewUserRepository(database database.Database) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (userRepository *UserRepository) Add(ctx context.Context, user models.User) error {

	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
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
		string(hashedPasswordBytes),
		user.CreatedAt,
		adminFlag,
	)
	return err
}

func (userRepository *UserRepository) Update(ctx context.Context, user models.User) error {

	if user.Password != nil {
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		_, err := userRepository.database.ExecContext(
			ctx,
			`
            UPDATE users SET
				email = ?,
				name = ?,
				password_hash = ?,
				updated_at = ?
			WHERE id = ?
        `,
			user.Email,
			user.Name,
			string(hashedPasswordBytes),
			user.UpdatedAt,
			user.ID,
		)
		return err
	} else {
		_, err := userRepository.database.ExecContext(
			ctx,
			`
            UPDATE users SET
				email = ?,
				name = ?,
				updated_at = ?
			WHERE id = ?
        `,
			user.Email,
			user.Name,
			user.UpdatedAt,
			user.ID,
		)
		return err
	}
}

func (userRepository *UserRepository) Delete(ctx context.Context, id string) error {
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

func (userRepository *UserRepository) Get(ctx context.Context, id string) (models.User, error) {
	var user models.User
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
	user.UpdatedAt = utils.Int64Ptr(updatedAt)
	user.IsSuperUser = isSuperUser.Valid && isSuperUser.Byte == 1
	user.AvatarURL = "https://i.pravatar.cc/48?id=" + user.ID
	return user, err
}

func (userRepository *UserRepository) Search(ctx context.Context) ([]models.User, error) {
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
	var users []models.User
	for rows.Next() {
		var user models.User
		var updatedAt sql.NullInt64
		var isSuperUser sql.NullByte
		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &updatedAt, &isSuperUser); err != nil {
			return nil, err
		}

		user.UpdatedAt = utils.Int64Ptr(updatedAt)
		user.IsSuperUser = isSuperUser.Valid && isSuperUser.Byte == 1
		user.AvatarURL = "https://i.pravatar.cc/48?id=" + user.ID
		users = append(users, user)
	}

	return users, nil
}
