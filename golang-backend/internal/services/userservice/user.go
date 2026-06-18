package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Add(ctx context.Context, user domain.User, password string) (domain.User, error)
	Update(ctx context.Context, user domain.User, password *string) (domain.User, error)
	Patch(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, userID string) error
	UnDelete(ctx context.Context, userID string) error
	Purge(ctx context.Context, userID string) error
	Get(ctx context.Context, userID string) (domain.User, error)
	SearchBase(ctx context.Context) ([]domain.UserBase, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error)
}

type userService struct {
	db                   database.Database
	permissionCache      cache.PermissionCache
	authorizationService authorizationservice.AuthorizationService
	userRepository       userrepository.UserRepository
}

func NewService(db database.Database, cache cache.PermissionCache, authorizationService authorizationservice.AuthorizationService, userRepository userrepository.UserRepository) UserService {
	return &userService{db: db, permissionCache: cache, authorizationService: authorizationService, userRepository: userRepository}
}

func (service *userService) Add(ctx context.Context, user domain.User, password string) (domain.User, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.User{}, err
	}
	user.ID = utils.UUID()
	user.CreatedAt = time.Now()
	if hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); hashErr != nil {
		return domain.User{}, hashErr
	} else {
		user.PasswordHash = string(hashedPasswordBytes)
		if err := service.userRepository.Add(ctx, service.db, user); err != nil {
			return domain.User{}, fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
		}
		return user, nil
	}
}

func (service *userService) Update(ctx context.Context, user domain.User, password *string) (domain.User, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.User{}, err
	}
	if password != nil {
		if hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost); hashErr != nil {
			return domain.User{}, hashErr
		} else {
			user.PasswordHash = string(hashedPasswordBytes)
		}
	}
	user.UpdatedAt = utils.NowToTimePtr()
	if err := service.userRepository.Update(ctx, service.db, user); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
	}
	cache.NewPermissionCache().DeleteUser(user.ID)
	return user, nil
}

func (service *userService) Patch(ctx context.Context, user domain.User) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if user.DeletedAt == nil {
		if err := service.userRepository.UnDelete(ctx, service.db, user.ID); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	} else {
		if err := service.userRepository.Delete(ctx, service.db, user.ID, time.Now().UnixMilli()); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	}
	cache.NewPermissionCache().DeleteUser(user.ID)
	return nil
}

func (service *userService) Delete(ctx context.Context, userID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.userRepository.Delete(ctx, service.db, userID, time.Now().UnixMilli()); err != nil {
		return err
	}
	cache.NewPermissionCache().DeleteUser(userID)
	return nil
}

func (service *userService) UnDelete(ctx context.Context, userID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.userRepository.UnDelete(ctx, service.db, userID); err != nil {
		return err
	}
	cache.NewPermissionCache().DeleteUser(userID)
	return nil
}

func (service *userService) Purge(ctx context.Context, userID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.userRepository.Purge(ctx, service.db, userID); err != nil {
		return err
	}
	cache.NewPermissionCache().DeleteUser(userID)
	return nil
}

func (service *userService) Get(ctx context.Context, userID string) (domain.User, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.User{}, err
	}
	if user, err := service.userRepository.Get(ctx, service.db, userID); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to get user with ID %s: %w", userID, err)
	} else {
		return user, nil
	}
}

func (service *userService) SearchBase(ctx context.Context) ([]domain.UserBase, error) {
	if users, err := service.userRepository.SearchBase(ctx, service.db); err != nil {
		return nil, fmt.Errorf("[UserService] failed to search base users: %w", err)
	} else {
		return users, nil
	}
}

func (service *userService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.Result{}, err
	}
	if users, pagerResult, err := service.userRepository.Search(ctx, service.db, pager, order, filter); err != nil {
		return nil, browser.Result{}, fmt.Errorf("[UserService] failed to search users: %w", err)
	} else {
		return users, pagerResult, nil
	}
}
