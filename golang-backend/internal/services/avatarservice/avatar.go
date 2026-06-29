package avatarservice

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/aportela/doneo/internal/middlewares"
	"github.com/disintegration/imaging"
)

type AvatarSize uint16

const (
	AvatarSizeTiny   AvatarSize = 32
	AvatarSizeSmall  AvatarSize = 64
	AvatarSizeNormal AvatarSize = 128
)

type AvatarService interface {
	GetUserAvatarPath(ctx context.Context, userID string, size AvatarSize) (string, error)
	SaveAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string) error
	SaveUserAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string, userID string, size AvatarSize) (string, error)
	DeleteAvatar(ctx context.Context) error
}

type avatarService struct {
	avatarBasePath string
}

func NewService(avatarBasePath string) AvatarService {
	return &avatarService{avatarBasePath: avatarBasePath}
}

func (service *avatarService) GetUserAvatarPath(ctx context.Context, userID string, size AvatarSize) (string, error) {
	if len(userID) != 32 {
		return "", fmt.Errorf("[AvatarService] invalid userID")
	}
	avatarFilename := userID + ".jpg"
	return path.Join(service.avatarBasePath, strconv.FormatUint(uint64(size), 10), avatarFilename), nil
}

func (service *avatarService) SaveAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string) error {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return fmt.Errorf("[AvatarService] user not found in context")
	}
	if data, err := io.ReadAll(sourceFile); err != nil {
		return fmt.Errorf("[AvatarService] cannot read file: %w", err)
	} else {
		if _, err := service.SaveUserAvatar(ctx, bytes.NewReader(data), sourceFilename, contextUser.ID, AvatarSizeTiny); err != nil {
			return fmt.Errorf("[AvatarService] error creating avatar tiny: %w", err)
		}
		if _, err := service.SaveUserAvatar(ctx, bytes.NewReader(data), sourceFilename, contextUser.ID, AvatarSizeSmall); err != nil {
			return fmt.Errorf("[AvatarService] error creating avatar small: %w", err)
		}
		if _, err := service.SaveUserAvatar(ctx, bytes.NewReader(data), sourceFilename, contextUser.ID, AvatarSizeNormal); err != nil {
			return fmt.Errorf("[AvatarService] error creating avatar normal: %w", err)
		}
	}
	return nil
}

func (service *avatarService) SaveUserAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string, userID string, size AvatarSize) (string, error) {
	avatarPath := path.Join(service.avatarBasePath, strconv.FormatUint(uint64(size), 10))
	if err := os.MkdirAll(avatarPath, 0755); err != nil {
		return "", err
	}
	fullPath := path.Join(avatarPath, userID+".jpg")
	if img, err := imaging.Decode(sourceFile); err != nil {
		return "", err
	} else {
		avatar := imaging.Fill(img, int(size), int(size), imaging.Center, imaging.Lanczos)
		if out, err := os.Create(fullPath); err != nil {
			return "", err
		} else {
			defer out.Close()
			if err := imaging.Encode(out, avatar, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
				return "", err
			}
			return fullPath, nil
		}
	}
}

func (service *avatarService) DeleteAvatar(ctx context.Context) error {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return fmt.Errorf("[AvatarService] user not found in context")
	}
	if err := os.MkdirAll(service.avatarBasePath, 0755); err != nil {
		return err
	}
	attachmentFilename := contextUser.ID + ".jpg"
	fullPath := path.Join(service.avatarBasePath, attachmentFilename)
	if _, err := os.Stat(fullPath); err != nil {
		return err
	}
	return os.Remove(fullPath)
}
