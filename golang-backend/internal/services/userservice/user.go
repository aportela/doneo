package userservice

import (
	"context"
	"database/sql"
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
	database             database.Database
	authorizationService authorizationservice.AuthorizationService
	repository           userrepository.UserRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, repository userrepository.UserRepository) UserService {
	return &userService{database: db, repository: repository}
}

func (service *userService) withUserAdminPermission(ctx context.Context, action func(userID string) error) error {
	currentUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[UserService] user not found in context")
	}

	if err := service.authorizationService.RequireUserPermission(ctx, currentUserID); err != nil {
		return err
	}

	return action(currentUserID)
}

func (service *userService) Add(ctx context.Context, user domain.User, password string) (domain.User, error) {
	// TODO: check user admin privilege ?
	user.ID = utils.UUID()
	user.CreatedAt = time.Now()
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}
	user.PasswordHash = string(hashedPasswordBytes)
	database.WithTx(ctx, func(tx *sql.Tx) error {
	})
	if err := service.repository.Add(ctx, service.database, user); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
	}
	return user, nil
}

func (service *userService) Update(ctx context.Context, user domain.User, password *string) (domain.User, error) {
	// TODO: check user admin privilege ?
	if password != nil {
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if hashErr != nil {
			return domain.User{}, hashErr
		}
		user.PasswordHash = string(hashedPasswordBytes)
	}
	user.UpdatedAt = utils.NowToTimePtr()
	if err := service.repository.Update(ctx, service.database, user); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
	}
	return user, nil
}

func (service *userService) Patch(ctx context.Context, user domain.User) error {
	// TODO: check user admin privilege ?
	if user.DeletedAt == nil {
		if err := service.repository.UnDelete(ctx, user.ID); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	} else {
		if err := service.repository.Delete(ctx, user.ID, time.Now().UnixMilli()); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	}
	return nil
}

func (service *userService) Delete(ctx context.Context, userID string) error {
	// TODO: check user admin privilege ?
	if err := service.repository.Delete(ctx, id, time.Now().UnixMilli()); err != nil {
		return fmt.Errorf("[UserService] failed to delete user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) UnDelete(ctx context.Context, userID string) error {
	// TODO: check user admin privilege ?
	if err := service.repository.UnDelete(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to undelete user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) Purge(ctx context.Context, userID string) error {
	// TODO: check user admin privilege ?
	if err := service.repository.Purge(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to purge user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) Get(ctx context.Context, userID string) (domain.User, error) {
	// TODO: check user admin privilege ?
	user, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to get user with ID %s: %w", id, err)
	}
	return user, nil
}

func (service *userService) SearchBase(ctx context.Context) ([]domain.UserBase, error) {

}

func (service *userService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	// TODO: check user admin privilege ?
	users, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[UserService] failed to search users: %w", err)
	}
	return users, pagerResult, nil
}
