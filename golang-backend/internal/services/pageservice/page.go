package pageservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/pagerepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type PageService interface {
	AddProjectPage(ctx context.Context, projectID string, page domain.Page) (domain.Page, error)
	UpdateProjectPage(ctx context.Context, projectID string, page domain.Page) (domain.Page, error)
	DeleteProjectPage(ctx context.Context, projectID string, pageID string) error
	GetProjectPage(ctx context.Context, projectID string, pageID string) (domain.Page, error)
	GetProjectPages(ctx context.Context, projectID string) ([]domain.Page, error)

	AddTaskPage(ctx context.Context, projectID string, taskID string, page domain.Page) (domain.Page, error)
	UpdateTaskPage(ctx context.Context, projectID string, taskID string, page domain.Page) (domain.Page, error)
	DeleteTaskPage(ctx context.Context, projectID string, taskID string, pageID string) error
	GetTaskPage(ctx context.Context, projectID string, taskID string, pageID string) (domain.Page, error)
	GetTaskPages(ctx context.Context, projectID string, taskID string) ([]domain.Page, error)
}

type pageService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	pageRepository          pagerepository.PageRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, repository pagerepository.PageRepository) PageService {
	return &pageService{db: db, historyOperationService: historyOperationService, authorizationService: authorizationService, pageRepository: repository}
}

func (service *pageService) AddProjectPage(ctx context.Context, projectID string, page domain.Page) (domain.Page, error) {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	} else {
		page.ID = utils.UUID()
		page.CreatedBy.ID = contextUser.ID
		page.CreatedBy.Name = contextUser.Name
		page.CreatedAt = time.Now()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.AddProjectPage(ctx, tx, projectID, page); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     page.CreatedAt,
					OperationType: domain.EventProjectPageAdded,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Page{}, err
		}
		return page, nil
	}
}

func (service *pageService) UpdateProjectPage(ctx context.Context, projectID string, page domain.Page) (domain.Page, error) {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	} else {
		page.UpdatedAt = utils.CurrentTimePtr()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.UpdateProjectPage(ctx, tx, projectID, page); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: contextUser.ID,
					},
					CreatedAt:     *page.UpdatedAt,
					OperationType: domain.EventProjectPageUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Page{}, err
		}
		return page, nil
	}
}

func (service *pageService) DeleteProjectPage(ctx context.Context, projectID string, pageID string) error {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.DeleteProjectPage(ctx, tx, projectID, pageID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: contextUser.ID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectPageDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *pageService) GetProjectPage(ctx context.Context, projectID string, pageID string) (domain.Page, error) {
	if _, err := service.authorizationService.RequireProjectViewPermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	}
	if notes, err := service.pageRepository.GetProjectPage(ctx, service.db, projectID, pageID); err != nil {
		return domain.Page{}, fmt.Errorf("[PageService] failed to get project page: %w", err)
	} else {
		return notes, nil
	}
}

func (service *pageService) GetProjectPages(ctx context.Context, projectID string) ([]domain.Page, error) {
	if _, err := service.authorizationService.RequireProjectViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if notes, err := service.pageRepository.GetProjectPages(ctx, service.db, projectID); err != nil {
		return nil, fmt.Errorf("[PageService] failed to get project pages: %w", err)
	} else {
		return notes, nil
	}
}

func (service *pageService) AddTaskPage(ctx context.Context, projectID string, taskID string, page domain.Page) (domain.Page, error) {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	} else {
		page.ID = utils.UUID()
		page.CreatedBy.ID = contextUser.ID
		page.CreatedBy.Name = contextUser.Name
		page.CreatedAt = time.Now()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.AddTaskPage(ctx, tx, taskID, page); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     page.CreatedAt,
					OperationType: domain.EventTaskPageAdded,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Page{}, err
		}
		return page, nil
	}
}

func (service *pageService) UpdateTaskPage(ctx context.Context, projectID string, taskID string, page domain.Page) (domain.Page, error) {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	} else {
		page.UpdatedAt = utils.CurrentTimePtr()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.UpdateTaskPage(ctx, tx, taskID, page); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: contextUser.ID,
					},
					CreatedAt:     *page.UpdatedAt,
					OperationType: domain.EventTaskPageUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Page{}, err
		}
		return page, nil
	}
}

func (service *pageService) DeleteTaskPage(ctx context.Context, projectID string, taskID string, pageID string) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.pageRepository.DeleteTaskPage(ctx, tx, taskID, pageID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: contextUser.ID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskPageDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *pageService) GetTaskPage(ctx context.Context, projectID string, taskID string, pageID string) (domain.Page, error) {
	if _, err := service.authorizationService.RequireTaskViewPermission(ctx, projectID); err != nil {
		return domain.Page{}, err
	}
	if notes, err := service.pageRepository.GetTaskPage(ctx, service.db, taskID, pageID); err != nil {
		return domain.Page{}, fmt.Errorf("[PageService] failed to get project page: %w", err)
	} else {
		return notes, nil
	}
}

func (service *pageService) GetTaskPages(ctx context.Context, projectID string, taskID string) ([]domain.Page, error) {
	if _, err := service.authorizationService.RequireProjectViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if notes, err := service.pageRepository.GetTaskPages(ctx, service.db, taskID); err != nil {
		return nil, fmt.Errorf("[PageService] failed to get task pages: %w", err)
	} else {
		return notes, nil
	}
}
