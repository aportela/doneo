package projecttypeservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectTypeService interface {
	Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error)
	Delete(ctx context.Context, projectTypeID string) error
	Get(ctx context.Context, projectTypeID string) (domain.ProjectType, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error)
}

type projectTypeService struct {
	db                    database.Database
	authorizationService  authorizationservice.AuthorizationService
	projectTypeRepository projecttyperepository.ProjectTypeRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, projectTypeRepository projecttyperepository.ProjectTypeRepository) ProjectTypeService {
	return &projectTypeService{db: db, authorizationService: authorizationService, projectTypeRepository: projectTypeRepository}
}

func (service *projectTypeService) withUserAdminPermission(ctx context.Context, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *projectTypeService) Add(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	err := service.withUserAdminPermission(ctx, func(currentUserID string) error {
		projectType.ID = utils.UUID()
		err := service.projectTypeRepository.Add(ctx, service.db, projectType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil
}

func (service *projectTypeService) Update(ctx context.Context, projectType domain.ProjectType) (domain.ProjectType, error) {
	err := service.withUserAdminPermission(ctx, func(currentUserID string) error {
		err := service.projectTypeRepository.Update(ctx, service.db, projectType)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return domain.ProjectType{}, err
	}
	return projectType, nil
}

func (service *projectTypeService) Delete(ctx context.Context, projectTypeID string) error {
	err := service.withUserAdminPermission(ctx, func(currentUserID string) error {
		err := service.projectTypeRepository.Delete(ctx, service.db, projectTypeID)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (service *projectTypeService) Get(ctx context.Context, projectTypeID string) (domain.ProjectType, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.ProjectType{}, fmt.Errorf("user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return domain.ProjectType{}, err
	}
	projectType, err := service.projectTypeRepository.Get(ctx, service.db, projectTypeID)
	if err != nil {
		return domain.ProjectType{}, fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", projectTypeID, err)
	}
	return projectType, nil
}

func (service *projectTypeService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectTypesFilter) ([]domain.ProjectType, browser.Result, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, browser.Result{}, fmt.Errorf("user not found in context")
	}

	if err := service.authorizationService.RequireUserAdminPermission(ctx, currentContextUserID); err != nil {
		return nil, browser.Result{}, err
	}
	projectTypes, pagerResult, err := service.projectTypeRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectTypeService] failed to search project types: %w", err)
	}
	return projectTypes, pagerResult, nil
}
