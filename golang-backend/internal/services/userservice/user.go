package userservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	AddUser(ctx context.Context, user domain.User) error
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, id string) error
	GetUser(ctx context.Context, id string) (domain.User, error)
	SearchUsers(ctx context.Context) ([]domain.User, error)
}

type userService struct {
	repository userrepository.UserRepository
}

func NewUserService(repository userrepository.UserRepository) UserService {
	return &userService{repository: repository}
}

func (s *userService) AddUser(ctx context.Context, user domain.User) error {
	hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hashErr != nil {
		return hashErr
	}
	user.PasswordHash = string(hashedPasswordBytes)
	user.CreatedAt = time.Now()
	if err := s.repository.Add(ctx, userrepository.UserToDTO(user)); err != nil {
		return fmt.Errorf("[UserService] failed to add user with ID %s: %w", user.ID, err)
	}
	return nil
}

func (s *userService) UpdateUser(ctx context.Context, user domain.User) error {
	if user.Password != "" {
		hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}
		user.PasswordHash = string(hashedPasswordBytes)
	}

	user.UpdatedAt = utils.NowToTimePtr()
	if err := s.repository.Update(ctx, userrepository.UserToDTO(user)); err != nil {
		return fmt.Errorf("[UserService] failed to update user with ID %s: %w", user.ID, err)
	}
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[UserService] failed to delete user with ID %s: %w", id, err)
	}
	return nil
}

func (s *userService) GetUser(ctx context.Context, id string) (domain.User, error) {
	user, err := s.repository.Get(ctx, id)
	if err != nil {
		return userrepository.DTOToUser(user), fmt.Errorf("[UserService] failed to get user with ID %s: %w", id, err)
	}
	return userrepository.DTOToUser(user), nil
}

func (s *userService) SearchUsers(ctx context.Context) ([]domain.User, error) {
	users, err := s.repository.Search(ctx)
	if err != nil {
		return nil, fmt.Errorf("[UserService] failed to search users: %w", err)
	}
	return userrepository.ToUserArray(users), nil
}
