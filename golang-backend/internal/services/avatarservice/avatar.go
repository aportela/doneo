package avatarservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/aportela/doneo/internal/middlewares"
)

const (
	DefaultAvatar string = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 128 128"><circle cx="64" cy="64" r="60" fill="#E5E7EB"/><circle cx="64" cy="46" r="18" fill="#6B7280"/><path d="M30 102C30 82 46 70 64 70s34 12 34 32Z" fill="#6B7280"/></svg>`
)

type AvatarService interface {
	GetUserAvatarPath(ctx context.Context, userID string) (string, error)
	GetContextUserAvatarPath(ctx context.Context) (string, error)
	SaveUserAvatar(ctx context.Context, svg string, userID string) error
	SaveContextUserAvatar(ctx context.Context, svg string) error
	DeleteUserAvatar(ctx context.Context, userID string) error
	DeleteContextUserAvatar(ctx context.Context) error
}

type avatarService struct {
	avatarBasePath string
}

func NewService(avatarBasePath string) AvatarService {
	return &avatarService{avatarBasePath: avatarBasePath}
}

func (service *avatarService) GetUserAvatarPath(ctx context.Context, userID string) (string, error) {
	if len(userID) != 36 {
		return "", fmt.Errorf("[AvatarService] invalid userID")
	}
	return path.Join(service.avatarBasePath, userID+".svg"), nil
}

func (service *avatarService) GetContextUserAvatarPath(ctx context.Context) (string, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return "", fmt.Errorf("[AvatarService] user not found in context")
	}
	return service.GetUserAvatarPath(ctx, contextUser.ID)
}

func (service *avatarService) SaveUserAvatar(ctx context.Context, svg string, userID string) error {
	if err := os.MkdirAll(service.avatarBasePath, 0755); err != nil {
		return err
	}
	avatarFullPath := path.Join(service.avatarBasePath, userID+".svg")

	if out, err := os.Create(avatarFullPath); err != nil {
		return err
	} else {
		defer out.Close()
		if _, err := out.WriteString(svg); err != nil {
			return err
		}
	}
	return nil
}

func (service *avatarService) SaveContextUserAvatar(ctx context.Context, svg string) error {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return fmt.Errorf("[AvatarService] user not found in context")
	}
	return service.SaveUserAvatar(ctx, svg, contextUser.ID)
}

func (service *avatarService) DeleteUserAvatar(ctx context.Context, userID string) error {
	avatarFullPath := path.Join(service.avatarBasePath, userID+".svg")
	if _, err := os.Stat(avatarFullPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}
	if err := os.Remove(avatarFullPath); err != nil {
		return err
	}
	return nil
}

func (service *avatarService) DeleteContextUserAvatar(ctx context.Context) error {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return fmt.Errorf("[AvatarService] user not found in context")
	}
	return service.DeleteUserAvatar(ctx, contextUser.ID)
}
