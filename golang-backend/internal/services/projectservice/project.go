package projectservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrepository"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectService interface {
	Add(ctx context.Context, Project domain.Project) (domain.Project, error)
	Update(ctx context.Context, Project domain.Project) (domain.Project, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Project, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error)
}

type projectService struct {
	database   database.Database
	repository projectrepository.ProjectRepository
}

func NewService(database database.Database, repository projectrepository.ProjectRepository) ProjectService {
	return &projectService{database: database, repository: repository}
}

func (service *projectService) Add(ctx context.Context, project domain.Project) (domain.Project, error) {
	tx, err := service.database.Begin()
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
	err = service.repository.Add(ctx, project)
	if err != nil {
		return domain.Project{}, err
	}
	err = service.repository.AddTaskCounter(ctx, project.ID)
	if err != nil {
		return domain.Project{}, err
	}
	err = projecthistoryrepository.NewRepository(service.database).AddProjectOperation(ctx, project.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: project.CreatedAt, OperationType: domain.EventProjectCreated})
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
	tx, err := service.database.Begin()
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
	project.UpdatedAt = utils.CurrentTimePtr()
	err = service.repository.Update(ctx, project)
	if err != nil {
		return domain.Project{}, err
	}
	err = projecthistoryrepository.NewRepository(service.database).AddProjectOperation(ctx, project.ID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectUpdated})
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
	tx, err := service.database.Begin()
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
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[ProjectService] user ID not found in context")
	}
	err = service.repository.Delete(ctx, id, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).AddProjectOperation(ctx, id, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *projectService) Get(ctx context.Context, id string) (domain.Project, error) {
	project, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.Project{}, fmt.Errorf("[ProjectService] failed to get project with ID %s: %w", id, err)
	}
	return project, nil
}

func (service *projectService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectFilter) ([]domain.Project, browser.Result, error) {
	projects, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectService] failed to search projects: %w", err)
	}
	return projects, pagerResult, nil
}
