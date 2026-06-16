package projectstatusservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectStatusService interface {
	Add(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error)
	Update(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error)
	Delete(ctx context.Context, projectStatusID string) error
	Get(ctx context.Context, projectStatusID string) (domain.ProjectStatus, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error)
}

type projectStatusService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	projectStatusRepository projectstatusrepository.ProjectStatusRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, projectStatusRepository projectstatusrepository.ProjectStatusRepository) ProjectStatusService {
	return &projectStatusService{db: db, authorizationService: authorizationService, projectStatusRepository: projectStatusRepository}
}

func (service *projectStatusService) Add(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectStatus{}, err
	}
	projectStatus.ID = utils.UUID()
	if err := service.projectStatusRepository.Add(ctx, service.db, projectStatus); err != nil {
		return domain.ProjectStatus{}, err
	}
	return projectStatus, nil
}

func (service *projectStatusService) Update(ctx context.Context, projectStatus domain.ProjectStatus) (domain.ProjectStatus, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectStatus{}, err
	}
	if err := service.projectStatusRepository.Update(ctx, service.db, projectStatus); err != nil {
		return domain.ProjectStatus{}, err
	}
	return projectStatus, nil
}

func (service *projectStatusService) Delete(ctx context.Context, projectStatusID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.projectStatusRepository.Delete(ctx, service.db, projectStatusID); err != nil {
		return err
	}
	return nil
}

func (service *projectStatusService) Get(ctx context.Context, projectStatusID string) (domain.ProjectStatus, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectStatus{}, err
	}
	projectStatus, err := service.projectStatusRepository.Get(ctx, service.db, projectStatusID)
	if err != nil {
		return domain.ProjectStatus{}, fmt.Errorf("[ProjectStatusService] failed to get project status with ID %s: %w", projectStatusID, err)
	}
	return projectStatus, nil
}

func (service *projectStatusService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectStatusesFilter) ([]domain.ProjectStatus, browser.Result, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.Result{}, err
	}
	projectStatuses, pagerResult, err := service.projectStatusRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectStatusService] failed to search project statuses: %w", err)
	}
	return projectStatuses, pagerResult, nil
}
