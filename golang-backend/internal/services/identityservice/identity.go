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
	GetCurrentUserInfo(ctx context.Context) (domain.User, error)
}

type identityService struct {
	db             database.Database
	userRepository userrepository.UserRepository
}

func NewService(db database.Database, repository userrepository.UserRepository) IdentityService {
	return &identityService{db: db, userRepository: repository}
}

func (service *identityService) SignIn(ctx context.Context, email string, password string) (domain.User, error) {
	user, err := service.userRepository.GetByEmail(ctx, service.db, email)
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
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[IdentityService] user not found in context")
	}
	user, err := service.userRepository.Get(ctx, service.db, contextUser.ID)
	if err != nil {
		return domain.User{}, err
	}
	if user.DeletedAt != nil {
		return domain.User{}, domain.DeletedError
	}
	return user, nil
}
