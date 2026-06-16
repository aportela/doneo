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
	if user, err := service.userRepository.GetByEmail(ctx, service.db, email); err != nil {
		return domain.User{}, err
	} else {
		if user.DeletedAt != nil {
			return domain.User{}, domain.DeletedError
		}
		if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			return domain.User{}, domain.InvalidCredentialsError
		}
		return user, nil
	}
}

func (service *identityService) GetCurrentUserInfo(ctx context.Context) (domain.User, error) {
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return domain.User{}, fmt.Errorf("[IdentityService] user not found in context")
	} else {
		if user, err := service.userRepository.Get(ctx, service.db, contextUser.ID); err != nil {
			return domain.User{}, err
		} else {
			if user.DeletedAt != nil {
				return domain.User{}, domain.DeletedError
			}
			return user, nil
		}
	}
}
