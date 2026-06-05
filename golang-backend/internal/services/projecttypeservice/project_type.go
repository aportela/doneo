package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
)

type ProjectTypeService interface {
	Add(ctx context.Context, projectType domain.ProjectType) error
	Update(ctx context.Context, projectType domain.ProjectType) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectType, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error)
}

type projectTypeService struct {
	database   database.Database
	repository projecttyperepository.ProjectTypeRepository
}

func NewService(database database.Database, repository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{database: database, repository: repository}
}

func (service *projectTypeService) Add(ctx context.Context, projectType domain.ProjectType) error {
	return service.repository.Add(ctx, projectType)
}

func (service *projectTypeService) Update(ctx context.Context, projectType domain.ProjectType) error {
	return service.repository.Update(ctx, projectType)
}

func (service *projectTypeService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *projectTypeService) Get(ctx context.Context, id string) (domain.ProjectType, error) {
	projectType, err := service.repository.Get(ctx, id)
	if err != nil {
		return projectType, fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", id, err)
	}
	return projectType, nil
}

func (service *projectTypeService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error) {
	projectTypes, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	}
	return projectTypes, pagerResult, nil
}
