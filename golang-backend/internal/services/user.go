package services

import (
	"context"

	"github.com/aportela/doneo/internal/models"
	"github.com/aportela/doneo/internal/repositories"
)

type UserService struct {
	repository *repositories.UserRepository
}

func NewUserService(repository *repositories.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) AddUser(ctx context.Context, user models.User) error {
	return s.repository.Add(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, user models.User) error {
	return s.repository.Update(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repository.Delete(ctx, id)
}

func (s *UserService) GetUser(ctx context.Context, id string) (models.User, error) {
	return s.repository.Get(ctx, id)
}

func (s *UserService) SearchUsers(ctx context.Context) ([]models.User, error) {
	return s.repository.Search(ctx)
}
