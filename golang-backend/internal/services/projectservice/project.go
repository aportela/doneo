package projectservice

import (
	"context"
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
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Project, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error)
}

type projectService struct {
	db      database.Database
	auth    authorizationservice.AuthorizationService
	history historyoperationservice.HistoryOperationService
	repo    projectrepository.ProjectRepository
}

func NewService(db database.Database, auth authorizationservice.AuthorizationService, history historyoperationservice.HistoryOperationService, repo projectrepository.ProjectRepository) ProjectService {
	return &projectService{db: db, auth: auth, history: history, repo: repo}
}

func (service *projectService) Add(ctx context.Context, project domain.Project) (domain.Project, error) {
	tx, err := service.db.Begin()
	if err != nil {
		return domain.Project{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Project{}, fmt.Errorf("[ProjectService] user ID not found in context")
	}
	project.ID = utils.UUID()
	project.CreatedBy.ID = currentUserId
	project.CreatedAt = time.Now()
	err = service.repo.Add(ctx, project)
	if err != nil {
		return domain.Project{}, err
	}
	err = service.repo.AddTaskCounter(ctx, project.ID)
	if err != nil {
		return domain.Project{}, err
	}
	_, err = service.history.AddProjectHistoryOperation(ctx, project.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: project.CreatedAt, OperationType: domain.EventProjectCreated})
	if err != nil {
		return domain.Project{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (service *projectService) Update(ctx context.Context, project domain.Project) (domain.Project, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Project{}, fmt.Errorf("[ProjectService] user ID not found in context")
	}
	err := service.auth.RequireProjectUpdatePermission(ctx, currentUserId, project.ID)
	if err != nil {
		return domain.Project{}, err
	}
	tx, err := service.db.Begin()
	if err != nil {
		return domain.Project{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	project.UpdatedAt = utils.CurrentTimePtr()
	err = service.repo.Update(ctx, project)
	if err != nil {
		return domain.Project{}, err
	}
	_, err = service.history.AddProjectHistoryOperation(ctx, project.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectUpdated})
	if err != nil {
		return domain.Project{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Project{}, err
	}
	return project, nil
}

func (service *projectService) Delete(ctx context.Context, id string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[ProjectService] user ID not found in context")
	}
	err := service.auth.RequireProjectDeletePermission(ctx, currentUserId, id)
	if err != nil {
		return err
	}
	tx, err := service.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	err = service.repo.Delete(ctx, id, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	_, err = service.history.AddProjectHistoryOperation(ctx, id, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *projectService) Get(ctx context.Context, id string) (domain.Project, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Project{}, fmt.Errorf("[ProjectService] user ID not found in context")
	}
	err := service.auth.RequireProjectViewPermission(ctx, currentUserId, id)
	if err != nil {
		return domain.Project{}, err
	}
	project, err := service.repo.Get(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("[ProjectService] failed to get project with ID %s: %w", id, err)
	}
	return project, nil
}

func (service *projectService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("[ProjectService] user ID not found in context")
	}
	if false {
		// TODO: user is not admin, search only projects with view permission for current user
		filter.ViewByUserId = &currentUserId
	}
	projects, pagerResult, err := service.repo.Search(ctx, pager, order, filter)

	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectService] failed to search projects: %w", err)
	}
	return projects, pagerResult, nil
}
