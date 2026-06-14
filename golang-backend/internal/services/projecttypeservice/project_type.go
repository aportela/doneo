package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectTypeService interface {
	Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.ProjectType, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error)
}

type projectTypeService struct {
	database   database.Database
	repository projecttyperepository.ProjectTypeRepository
}

func NewService(db database.Database, repository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{database: db, repository: repository}
}

func (service *projectTypeService) Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	projectType.ID = utils.UUID()
	err := service.repository.Add(ctx, projectType)
	if err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil
}

func (service *projectTypeService) Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	err := service.repository.Update(ctx, projectType)
	if err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil
}

func (service *projectTypeService) Delete(ctx context.Context, id string) error {
	return service.repository.Delete(ctx, id)
}

func (service *projectTypeService) Get(ctx context.Context, id string) (domain.ProjectType, error) {
	projectType, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.ProjectType{}, fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", id, err)
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
