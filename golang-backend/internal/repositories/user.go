package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(database *sql.DB) *UserRepository {
	return &UserRepository{
		database: database,
	}
}

func (userRepository *UserRepository) Add(ctx context.Context, user models.User) error {

	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO USER (id, email, name, password_hash, ctime, mtime)
			VALUES (?, ?, ?, ?, ?, NULL)
        `,
		user.ID,
		user.Email,
		user.Name,
		hashedPassword,
		utils.CurrentTimestamp(),
	)
	return err
}

func (userRepository *UserRepository) Update(ctx context.Context, user models.User) error {

	if user.Password != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		_, err := userRepository.database.ExecContext(
			ctx,
			`
            UPDATE USER SET
				email = ?,
				name = ?,
				password_hash = ?,
				mtime = ?
			WHERE id = ?
        `,
			user.Email,
			user.Name,
			hashedPassword,
			utils.CurrentTimestamp(),
			user.ID,
		)

		return err
	} else {
		_, err := userRepository.database.ExecContext(
			ctx,
			`
            UPDATE USER SET
				email = ?,
				name = ?,
				mtime = ?
			WHERE id = ?
        `,
			user.Email,
			user.Name,
			utils.CurrentTimestamp(),
			user.ID,
		)

		return err
	}
}

func (userRepository *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := userRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM USER
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (userRepository *UserRepository) Get(ctx context.Context, id string) (models.User, error) {
	var user models.User
	var mtime sql.NullInt64
	err := userRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                U.id, U.email, U.name, U.ctime, U.mtime
            FROM USER U
            WHERE U.id = ?
        `,
		id).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &mtime)
	user.LastUpdateAt = utils.Int64Ptr(mtime)

	return user, err
}

func (userRepository *UserRepository) Search(ctx context.Context) ([]models.User, error) {
	rows, err := userRepository.database.QueryContext(
		ctx,
		`
			SELECT
				U.id, U.email, U.name, U.ctime, U.mtime
			FROM USER U
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		var mtime sql.NullInt64

		if err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &mtime); err != nil {
			return nil, err
		}

		user.LastUpdateAt = utils.Int64Ptr(mtime)
		users = append(users, user)
	}

	return users, nil
}
