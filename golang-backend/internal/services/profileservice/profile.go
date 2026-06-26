package profileservice

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/disintegration/imaging"
	"golang.org/x/crypto/bcrypt"
)

type ProfileService interface {
	Update(ctx context.Context, user domain.User, password *string) (domain.User, error)
	Get(ctx context.Context) (domain.User, error)
	SaveAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string) (string, error)
	DeleteAvatar(ctx context.Context) error
}

type profileService struct {
	db                   database.Database
	avatarBasePath       string
	permissionCache      cache.PermissionCache
	authorizationService authorizationservice.AuthorizationService
	userRepository       userrepository.UserRepository
}

func NewService(db database.Database, avatarBasePath string, cache cache.PermissionCache, authorizationService authorizationservice.AuthorizationService, userRepository userrepository.UserRepository) ProfileService {
	return &profileService{db: db, avatarBasePath: avatarBasePath, permissionCache: cache, authorizationService: authorizationService, userRepository: userRepository}
}

func (service *profileService) Update(ctx context.Context, user domain.User, password *string) (domain.User, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[ProfileService] user not found in context")
	}
	if currentUser, err := service.userRepository.Get(ctx, service.db, contextUser.ID); err != nil {
		return domain.User{}, fmt.Errorf("[ProfileService] error getting current user data")
	} else {
		if password != nil {
			if hashedPasswordBytes, hashErr := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost); hashErr != nil {
				return domain.User{}, hashErr
			} else {
				currentUser.PasswordHash = string(hashedPasswordBytes)
			}
		}
		currentUser.Name = user.Name
		currentUser.Email = user.Email
		currentUser.UpdatedAt = utils.NowToTimePtr()
		if err := service.userRepository.Update(ctx, service.db, currentUser); err != nil {
			return domain.User{}, fmt.Errorf("[ProfileService] failed to update user with ID %s: %w", user.ID, err)
		}
		cache.NewPermissionCache().DeleteUser(currentUser.ID)
		return currentUser, nil
	}
}

func (service *profileService) Get(ctx context.Context) (domain.User, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return domain.User{}, fmt.Errorf("[ProfileService] user not found in context")
	}
	if user, err := service.userRepository.Get(ctx, service.db, contextUser.ID); err != nil {
		return domain.User{}, fmt.Errorf("[ProfileService] failed to get user with ID %s: %w", contextUser.ID, err)
	} else {
		return user, nil
	}
}

func (service *profileService) SaveAvatar(ctx context.Context, sourceFile io.Reader, sourceFilename string) (string, error) {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return "", fmt.Errorf("[ProfileService] user not found in context")
	}
	if err := os.MkdirAll(service.avatarBasePath, 0755); err != nil {
		return "", err
	}
	attachmentFilename := contextUser.ID + ".jpg"
	fullPath := path.Join(service.avatarBasePath, attachmentFilename)
	if img, err := imaging.Decode(sourceFile); err != nil {
		return "", err
	} else {
		avatar := imaging.Fill(img, 100, 100, imaging.Center, imaging.Lanczos)
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

func (service *profileService) DeleteAvatar(ctx context.Context) error {
	contextUser, ok := middlewares.GetContextUser(ctx)
	if !ok {
		return fmt.Errorf("[ProfileService] user not found in context")
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
