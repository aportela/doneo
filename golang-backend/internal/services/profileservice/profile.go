package profileservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	Update(ctx context.Context, user domain.User, password *string) (domain.User, error)
	Get(ctx context.Context) (domain.User, error)
}

type profileService struct {
	db                   database.Database
	permissionCache      cache.PermissionCache
	authorizationService authorizationservice.AuthorizationService
	userRepository       userrepository.UserRepository
}

func NewService(db database.Database, cache cache.PermissionCache, authorizationService authorizationservice.AuthorizationService, userRepository userrepository.UserRepository) ProfileService {
	return &profileService{db: db, permissionCache: cache, authorizationService: authorizationService, userRepository: userRepository}
}

func (service *profileService) Update(ctx context.Context, user domain.User, password *string) (domain.User, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[ProfileService] user not found in context")
	}
	if currentUser, err := service.userRepository.Get(ctx, service.db, contextUser.ID); err != nil {
		return domain.User{}, fmt.Errorf("[ProfileService] error getting current user data")
	} else {
		if password != nil {
			if hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost); hashErr != nil {
				return domain.User{}, hashErr
			} else {
				currentUser.PasswordHash = string(hashedPasswordBytes)
			}
		}
		currentUser.Name = user.Name
		currentUser.Email = user.Email
		currentUser.UpdatedAt = utils.NowToTimePtr()
		if err := service.userRepository.Update(ctx, service.db, currentUser); err != nil {
			return domain.User{}, fmt.Errorf("[ProfileService] failed to update user with ID %s: %w", user.ID, err)
		}
		cache.NewPermissionCache().DeleteUser(currentUser.ID)
		return currentUser, nil
	}
}

func (service *profileService) Get(ctx context.Context) (domain.User, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[ProfileService] user not found in context")
	}
	if user, err := service.userRepository.Get(ctx, service.db, contextUser.ID); err != nil {
		return domain.User{}, fmt.Errorf("[ProfileService] failed to get user with ID %s: %w", contextUser.ID, err)
	} else {
		return user, nil
	}
}
