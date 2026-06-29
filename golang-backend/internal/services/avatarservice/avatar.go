package avatarservice

import (
	"context"
	"fmt"
	"path"
	"strconv"
)

type AvatarSize uint16

const (
	AvatarSizeSmall  AvatarSize = 50
	AvatarSizeNormal AvatarSize = 100
	AvatarSizeLarge  AvatarSize = 200
)

type AvatarService interface {
	getAvatarFilename(userID string) string
	getAvatarPath(userID string, size AvatarSize) string
	GetUserAvatarPath(ctx context.Context, userID string, size AvatarSize) (string, error)
}

type avatarService struct {
	avatarBasePath string
}

func NewService(avatarBasePath string) AvatarService {
	return &avatarService{avatarBasePath: avatarBasePath}
}

func (service *avatarService) getAvatarFilename(userID string) string {
	return userID + ".jpg"
}

func (service *avatarService) getAvatarPath(userID string, size AvatarSize) string {
	return path.Join(service.avatarBasePath, strconv.FormatUint(uint64(size), 10), service.getAvatarFilename(userID))
}

func (service *avatarService) GetUserAvatarPath(ctx context.Context, userID string, size AvatarSize) (string, error) {
	if len(userID) != 32 {
		return "", fmt.Errorf("[AvatarService] invalid userID")
	}
	return service.getAvatarPath(userID, size), nil
}
