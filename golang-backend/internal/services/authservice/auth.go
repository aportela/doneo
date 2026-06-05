package authservice

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(ctx context.Context, email string, password string) (domain.User, error)
	GetUserInfo(ctx context.Context, userId string) (domain.User, error)
}

type authService struct {
	database   database.Database
	repository userrepository.UserRepository
}

func NewService(database database.Database, repository userrepository.UserRepository) AuthService {
	return &authService{database: database, repository: repository}
}

func (service *authService) SignIn(ctx context.Context, email string, password string) (domain.User, error) {
	// TODO: remove PasswordHash from domain & return userId, email, password Hash ?
	user, err := service.repository.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	if user.DeletedAt != nil {
		return domain.User{}, domain.DeletedError
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return domain.User{}, domain.InvalidCredentialsError
	}
	return user, nil
}

func (service *authService) GetUserInfo(ctx context.Context, userId string) (domain.User, error) {
	user, err := service.repository.Get(ctx, userId)
	if err != nil {
		return domain.User{}, err
	}
	if user.DeletedAt != nil {
		return domain.User{}, domain.DeletedError
	}
	return user, nil
}
