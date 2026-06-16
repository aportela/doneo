package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
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
	authorizationService authorizationservice.AuthorizationService
	userRepository       userrepository.UserRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, userRepository userrepository.UserRepository) UserService {
	return &userService{db: db, authorizationService: authorizationService, userRepository: userRepository}
}

func (service *userService) Add(ctx context.Context, user domain.User, password string) (domain.User, error) {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		user.ID = utils.UUID()
		user.CreatedAt = time.Now()
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		user.PasswordHash = string(hashedPasswordBytes)
		if err := service.userRepository.Add(ctx, service.db, user); err != nil {
			return fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
		}
		return nil
	})
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (service *userService) Update(ctx context.Context, user domain.User, password *string) (domain.User, error) {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		if password != nil {
			hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
			if hashErr != nil {
				return hashErr
			}
			user.PasswordHash = string(hashedPasswordBytes)
		}
		user.UpdatedAt = utils.NowToTimePtr()
		if err := service.userRepository.Update(ctx, service.db, user); err != nil {
			return fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
		}
		return nil
	})
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (service *userService) Patch(ctx context.Context, user domain.User) error {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		if user.DeletedAt == nil {
			if err := service.userRepository.UnDelete(ctx, service.db, user.ID); err != nil {
				return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
			}
		} else {
			if err := service.userRepository.Delete(ctx, service.db, user.ID, time.Now().UnixMilli()); err != nil {
				return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
			}
		}
		return nil
	})
	return err
}

func (service *userService) Delete(ctx context.Context, userID string) error {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		if err := service.userRepository.Delete(ctx, service.db, userID, time.Now().UnixMilli()); err != nil {
			return fmt.Errorf("[UserService] failed to delete user with ID %s: %w", userID, err)
		}
		return nil
	})
	return err
}

func (service *userService) UnDelete(ctx context.Context, userID string) error {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		if err := service.userRepository.UnDelete(ctx, service.db, userID); err != nil {
			return fmt.Errorf("[UserService] failed to undelete user with ID %s: %w", userID, err)
		}
		return nil
	})
	return err
}

func (service *userService) Purge(ctx context.Context, userID string) error {
	err := service.authorizationService.WithUserAdminPermission(ctx, func(currentUserID string) error {
		if err := service.userRepository.Purge(ctx, service.db, userID); err != nil {
			return fmt.Errorf("[UserService] failed to purge user with ID %s: %w", userID, err)
		}
		return nil
	})
	return err
}

func (service *userService) Get(ctx context.Context, userID string) (domain.User, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[UserService] user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return domain.User{}, err
	}
	user, err := service.userRepository.Get(ctx, service.db, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to get user with ID %s: %w", userID, err)
	}
	return user, nil
}

func (service *userService) SearchBase(ctx context.Context) ([]domain.UserBase, error) {
	users, err := service.userRepository.SearchBase(ctx, service.db)
	if err != nil {
		return nil, fmt.Errorf("[UserService] failed to search base users: %w", err)
	}
	return users, nil
}

func (service *userService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("[UserService] user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return nil, browser.Result{}, err
	}
	users, pagerResult, err := service.userRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[UserService] failed to search users: %w", err)
	}
	return users, pagerResult, nil
}
