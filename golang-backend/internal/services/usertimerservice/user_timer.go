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
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	} else {
		return service.userTimerRepository.StartUserTimer(ctx, service.db, utils.UUID(), contextUser.ID, summary, time.Now().UnixMilli())
	}
}

func (service *userTimerService) StopUserTimer(ctx context.Context, userTimerID string) error {
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	} else {
		return service.userTimerRepository.StopUserTimer(ctx, service.db, userTimerID, contextUser.ID, time.Now().UnixMilli())
	}
}

func (service *userTimerService) DeleteUserTimer(ctx context.Context, userTimerID string) error {
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	} else {
		return service.userTimerRepository.DeleteUserTimer(ctx, service.db, userTimerID, contextUser.ID)
	}
}

func (service *userTimerService) ClearUserTimers(ctx context.Context) error {
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return fmt.Errorf("[UserTimerService] user not found in context")
	} else {
		return service.userTimerRepository.ClearUserTimers(ctx, service.db, contextUser.ID)
	}
}

func (service *userTimerService) GetUserTimers(ctx context.Context) ([]domain.UserTimer, error) {
	if contextUser, ok := middlewares.GetContextUser(ctx); !ok {
		return nil, fmt.Errorf("[UserTimerService] user not found in context")
	} else {
		if userTimers, err := service.userTimerRepository.GetUserTimers(ctx, service.db, contextUser.ID); err != nil {
			return nil, fmt.Errorf("[UserTimerService] failed to get user timers: %w", err)
		} else {
			return userTimers, nil
		}
	}
}
