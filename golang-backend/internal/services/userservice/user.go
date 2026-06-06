package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Add(ctx context.Context, user domain.User, password string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	Patch(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	UnDelete(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.User, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error)
}

type userService struct {
	database   database.Database
	repository userrepository.UserRepository
}

func NewService(database database.Database, repository userrepository.UserRepository) UserService {
	return &userService{database: database, repository: repository}
}

func (service *userService) Add(ctx context.Context, user domain.User, password string) (domain.User, error) {
	user.ID = utils.UUID()
	user.CreatedAt = time.Now()
	if err := service.repository.Add(ctx, user, password); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
	}
	return user, nil
}

func (service *userService) Update(ctx context.Context, user domain.User) (domain.User, error) {
	if user.Password != "" {
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return domain.User{}, hashErr
		}
		user.PasswordHash = string(hashedPasswordBytes)
	}
	user.UpdatedAt = utils.NowToTimePtr()
	if err := service.repository.Update(ctx, user); err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
	}
	return user, nil
}

func (service *userService) Patch(ctx context.Context, user domain.User) error {
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

func (service *userService) Delete(ctx context.Context, id string) error {
	if err := service.repository.Delete(ctx, id, time.Now().UnixMilli()); err != nil {
		return fmt.Errorf("[UserService] failed to delete user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) UnDelete(ctx context.Context, id string) error {
	if err := service.repository.UnDelete(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to undelete user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) Purge(ctx context.Context, id string) error {
	if err := service.repository.Purge(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to purge user with ID %s: %w", id, err)
	}
	return nil
}

func (service *userService) Get(ctx context.Context, id string) (domain.User, error) {
	user, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to get user with ID %s: %w", id, err)
	}
	return user, nil
}

func (service *userService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	users, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[UserService] failed to search users: %w", err)
	}
	return users, pagerResult, nil
}
