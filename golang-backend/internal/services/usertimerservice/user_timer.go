package usertimerservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/usertimerrepository"
	"github.com/aportela/doneo/internal/utils"
)

type UserTimerService interface {
	StartUserTimer(ctx context.Context, summary string) error
	StopUserTimer(ctx context.Context, userTimerID string) error
	DeleteUserTimer(ctx context.Context, userTimerID string) error
	ClearUserTimers(ctx context.Context) error
	GetUserTimers(ctx context.Context) ([]domain.UserTimer, error)
}

type userTimerService struct {
	db                  database.Database
	userTimerRepository usertimerrepository.UserTimerRepository
}

func NewService(db database.Database, userTimerRepository usertimerrepository.UserTimerRepository) UserTimerService {
	return &userTimerService{db: db, userTimerRepository: userTimerRepository}
}

func (service *userTimerService) StartUserTimer(ctx context.Context, summary string) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	}
	err := service.userTimerRepository.StartUserTimer(ctx, utils.UUID(), currentContextUserID, summary, time.Now().UnixMilli())
	return err
}

func (service *userTimerService) StopUserTimer(ctx context.Context, userTimerID string) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	}
	err := service.userTimerRepository.StopUserTimer(ctx, userTimerID, currentContextUserID, time.Now().UnixMilli())
	return err
}

func (service *userTimerService) DeleteUserTimer(ctx context.Context, userTimerID string) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	}
	err := service.userTimerRepository.DeleteUserTimer(ctx, userTimerID, currentContextUserID)
	return err
}

func (service *userTimerService) ClearUserTimers(ctx context.Context) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	}
	err := service.userTimerRepository.ClearUserTimers(ctx, service.db, currentContextUserID)
	return err
}

func (service *userTimerService) GetUserTimers(ctx context.Context) ([]domain.UserTimer, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[UserTimerService] user not found in context")
	}
	userTimers, err := service.userTimerRepository.GetUserTimers(ctx, service.db, currentContextUserID)
	if err != nil {
		return nil, fmt.Errorf("[UserTimerService] failed to get user timers: %w", err)
	}
	return userTimers, nil
}
