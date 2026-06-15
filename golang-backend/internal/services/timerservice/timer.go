package timerservice

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

type TimerService interface {
	Start(ctx context.Context, summary string) error
	Stop(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Clear(ctx context.Context) error
	Search(ctx context.Context) ([]domain.UserTimer, error)
}

type timerService struct {
	database   database.Database
	repository timerrepository.UserTimerRepository
}

func NewService(db database.Database, repository timerrepository.UserTimerRepository) TimerService {
	return &timerService{database: db, repository: repository}
}

func (service *timerService) Start(ctx context.Context, summary string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.repository.Start(ctx, utils.UUID(), currentUserId, summary, time.Now().UnixMilli())
	return err
}

func (service *timerService) Stop(ctx context.Context, id string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.repository.Stop(ctx, id, currentUserId, time.Now().UnixMilli())
	return err
}

func (service *timerService) Delete(ctx context.Context, id string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.repository.Delete(ctx, id, currentUserId)
	return err
}

func (service *timerService) Clear(ctx context.Context) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[TimerService] user ID not found in context")
	}
	err := service.repository.Clear(ctx, currentUserId)
	return err
}

func (service *timerService) Search(ctx context.Context) ([]domain.UserTimer, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[TimerService] user ID not found in context")
	}
	timers, err := service.repository.Search(ctx, currentUserId)
	if err != nil {
		return nil, fmt.Errorf("[TimerService] failed to get user timers: %w", err)
	}
	return timers, nil
}
