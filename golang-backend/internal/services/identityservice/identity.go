package identityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"golang.org/x/crypto/bcrypt"
)

type IdentityService interface {
	SignIn(ctx context.Context, email string, password string) (domain.User, error)
	GetCurrentUserInfo(ctx context.Context, userId string) (domain.User, error)
}

type identityService struct {
	db         database.Database
	repository userrepository.UserRepository
}

func NewService(db database.Database, repository userrepository.UserRepository) IdentityService {
	return &identityService{db: db, repository: repository}
}

func (service *identityService) SignIn(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := service.repository.GetByEmail(ctx, service.db, email)
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

func (service *identityService) GetCurrentUserInfo(ctx context.Context) (domain.User, error) {
	currentUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("user not found in context")
	}
	user, err := service.repository.Get(ctx, service.db, currentUserID)
	if err != nil {
		return domain.User{}, err
	}
	if user.DeletedAt != nil {
		return domain.User{}, domain.DeletedError
	}
	return user, nil
}
