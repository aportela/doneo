package usertimerservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/timerrepository"
	"github.com/aportela/doneo/internal/utils"
)

type UserTimerService interface {
	Start(ctx context.Context, summary string) error
	Stop(ctx context.Context, userTimerID string) error
	Delete(ctx context.Context, userTimerID string) error
	ClearUserTimers(ctx context.Context) error
	GetUserTimers(ctx context.Context) ([]domain.UserTimer, error)
}

type userTimerService struct {
	db                  database.Database
	userTimerRepository timerrepository.UserTimerRepository
}

func NewService(db database.Database, userTimerRepository timerrepository.UserTimerRepository) UserTimerService {
	return &userTimerService{db: db, userTimerRepository: userTimerRepository}
}

func (service *userTimerService) Start(ctx context.Context, summary string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.userTimerRepository.Start(ctx, utils.UUID(), currentUserId, summary, time.Now().UnixMilli())
	return err
}

func (service *userTimerService) Stop(ctx context.Context, id string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.userTimerRepository.Stop(ctx, id, currentUserId, time.Now().UnixMilli())
	return err
}

func (service *userTimerService) Delete(ctx context.Context, id string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.userTimerRepository.Delete(ctx, id, currentUserId)
	return err
}

func (service *userTimerService) Clear(ctx context.Context) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.userTimerRepository.Clear(ctx, currentUserId)
	return err
}

func (service *userTimerService) Search(ctx context.Context) ([]domain.UserTimer, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[TimerService] user ID not found in context")
	}
	timers, err := service.userTimerRepository.Search(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("[TimerService] failed to get user timers: %w", err)
	}
	return timers, nil
}
