package authservice

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	SignIn(ctx context.Context, user domain.User) (domain.User, error)
	GetUserInfo(ctx context.Context, userId string) (domain.User, error)
}

type authService struct {
	database   database.Database
	repository userrepository.UserRepository
}

func NewService(db database.Database, repository userrepository.UserRepository) AuthService {
	return &authService{database: db, repository: repository}
}

func (service *authService) SignIn(ctx context.Context, user domain.User) (domain.User, error) {
	// TODO: remove PasswordHash from domain & return userId, email, password Hash ?
	credentialUser, err := service.repository.GetByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	if user.DeletedAt != nil {
		return domain.User{}, domain.DeletedError
	}
	err = bcrypt.CompareHashAndPassword([]byte(credentialUser.PasswordHash), []byte(user.Password))
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
