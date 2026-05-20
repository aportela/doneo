package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Add(ctx context.Context, user domain.User) error
	Update(ctx context.Context, user domain.User) error
	Patch(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	UnDelete(ctx context.Context, id string) error
	Purge(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.User, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error)
}

type userService struct {
	repository userrepository.UserRepository
}

func NewUserService(repository userrepository.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) Add(ctx context.Context, user domain.User) error {
	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	user.PasswordHash = string(hashedPasswordBytes)
	user.CreatedAt = time.Now()
	if err := s.repository.Add(ctx, userrepository.DomainToDTO(user)); err != nil {
		return fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
	}
	return nil
}

func (s *userService) Update(ctx context.Context, user domain.User) error {
	if user.Password != "" {
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		user.PasswordHash = string(hashedPasswordBytes)
	}
	user.UpdatedAt = utils.NowToTimePtr()
	if err := s.repository.Update(ctx, userrepository.DomainToDTO(user)); err != nil {
		return fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
	}
	return nil
}

func (s *userService) Patch(ctx context.Context, user domain.User) error {
	if user.DeletedAt == nil {
		if err := s.repository.UnDelete(ctx, user.ID); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	} else {
		if err := s.repository.Delete(ctx, user.ID); err != nil {
			return fmt.Errorf("[UserService] failed to patch user with ID %s: %w", user.ID, err)
		}
	}
	return nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to delete user with ID %s: %w", id, err)
	}
	return nil
}

func (s *userService) UnDelete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to undelete user with ID %s: %w", id, err)
	}
	return nil
}

func (s *userService) Purge(ctx context.Context, id string) error {
	if err := s.repository.Purge(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to purge user with ID %s: %w", id, err)
	}
	return nil
}

func (s *userService) Get(ctx context.Context, id string) (domain.User, error) {
	user, err := s.repository.Get(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("[UserService] failed to get user with ID %s: %w", id, err)
	}
	return userrepository.DTOToDomain(user), nil
}

func (s *userService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchUsersFilter) ([]domain.User, browser.Result, error) {
	users, pagerResult, err := s.repository.Search(ctx, pager, order, userrepository.DomainFilterToDTO(filter))
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[UserService] failed to search users: %w", err)
	}
	return userrepository.DTOArrayToDomainArray(users), pagerResult, nil
}
