package projectservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectService interface {
	Add(ctx context.Context, project domain.Project) (domain.Project, error)
	Update(ctx context.Context, project domain.Project) (domain.Project, error)
	Delete(ctx context.Context, projectID string) error
	Get(ctx context.Context, projectID string) (domain.Project, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error)
}

type projectService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	projectRepository       projectrepository.ProjectRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, projectRepository projectrepository.ProjectRepository) ProjectService {
	return &projectService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, projectRepository: projectRepository}
}

func (service *projectService) Add(ctx context.Context, project domain.Project) (domain.Project, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Project{}, fmt.Errorf("user not found in context")
	}
	project.ID = utils.UUID()
	project.CreatedBy.ID = currentContextUserID
	project.CreatedAt = time.Now()

	err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
		if err := service.projectRepository.Add(ctx, tx, project); err != nil {
			return err
		}
		if err := service.projectRepository.AddTaskCounter(ctx, tx, project.ID); err != nil {
			return err
		}
		if _, err := service.historyOperationService.AddProjectHistoryOperation(
			ctx,
			tx,
			project.ID,
			domain.HistoryOperation{
				ID:            utils.UUID(),
				CreatedBy:     domain.UserBase{ID: currentContextUserID},
				CreatedAt:     project.CreatedAt,
				OperationType: domain.EventProjectCreated,
			},
		); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (service *projectService) Update(ctx context.Context, project domain.Project) (domain.Project, error) {
	err := service.authorizationService.WithProjectUpdatePermission(ctx, project.ID, func(currentUserID string) error {
		project.UpdatedAt = utils.CurrentTimePtr()

		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.projectRepository.Update(ctx, tx, project); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				project.ID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     *project.UpdatedAt,
					OperationType: domain.EventProjectUpdated,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})
	if err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (service *projectService) Delete(ctx context.Context, projectID string) error {
	err := service.authorizationService.WithProjectDeletePermission(ctx, projectID, func(currentUserID string) error {

		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			deletedAt := time.Now()
			if err := service.projectRepository.Delete(ctx, tx, projectID, deletedAt.UnixMilli()); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     deletedAt,
					OperationType: domain.EventProjectDeleted,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})
	return err
}

func (service *projectService) Get(ctx context.Context, projectID string) (domain.Project, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Project{}, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return domain.Project{}, err
	}
	project, err := service.projectRepository.Get(ctx, service.db, projectID)
	if err != nil {
		return domain.Project{}, fmt.Errorf("[ProjectService] failed to get project with ID %s: %w", projectID, err)
	}
	return project, nil
}

func (service *projectService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		// filter by projects visible by current user when admin flag is not set
		filter.ViewByUserId = &currentContextUserID
	}
	projects, pagerResult, err := service.projectRepository.Search(ctx, service.db, pager, order, filter)

	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectService] failed to search projects: %w", err)
	}
	return projects, pagerResult, nil
}
