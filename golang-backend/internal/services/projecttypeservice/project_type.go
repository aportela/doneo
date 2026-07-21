package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectTypeService interface {
	Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Delete(ctx context.Context, projectTypeID string) error
	Get(ctx context.Context, projectTypeID string) (domain.ProjectType, error)
	SearchBase(ctx context.Context) ([]domain.ProjectType, browser.PagerResult, error)
	Search(ctx context.Context, pager browser.PagerQuery, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.PagerResult, error)
}

type projectTypeService struct {
	db                    database.Database
	authorizationService  authorizationservice.AuthorizationService
	projectTypeRepository projecttyperepository.ProjectTypeRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, projectTypeRepository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{db: db, authorizationService: authorizationService, projectTypeRepository: projectTypeRepository}
}

func (service *projectTypeService) Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectType{}, err
	}
	projectType.ID = utils.UUID()
	if err := service.projectTypeRepository.Add(ctx, service.db, projectType); err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil

}

func (service *projectTypeService) Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectType{}, err
	}
	if err := service.projectTypeRepository.Update(ctx, service.db, projectType); err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil
}

func (service *projectTypeService) Delete(ctx context.Context, projectTypeID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	return service.projectTypeRepository.Delete(ctx, service.db, projectTypeID)
}

func (service *projectTypeService) Get(ctx context.Context, projectTypeID string) (domain.ProjectType, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectType{}, err
	}
	if projectType, err := service.projectTypeRepository.Get(ctx, service.db, projectTypeID); err != nil {
		return domain.ProjectType{}, fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", projectTypeID, err)
	} else {
		return projectType, nil
	}
}

func (service *projectTypeService) SearchBase(ctx context.Context) ([]domain.ProjectType, browser.PagerResult, error) {
	if projectTypes, pagerResult, err := service.projectTypeRepository.Search(ctx, service.db, browser.PagerQuery{Enabled: false, CurrentPage: 1, ResultsPage: 0}, browser.Order{Field: "index", Direction: "ASC"}, domain.SearchProjectTypesFilter{}); err != nil {
		return nil, browser.PagerResult{}, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	} else {
		return projectTypes, pagerResult, nil
	}
}

func (service *projectTypeService) Search(ctx context.Context, pager browser.PagerQuery, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.PagerResult, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.PagerResult{}, err
	}
	if projectTypes, pagerResult, err := service.projectTypeRepository.Search(ctx, service.db, pager, order, filter); err != nil {
		return nil, browser.PagerResult{}, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	} else {
		return projectTypes, pagerResult, nil
	}
}
