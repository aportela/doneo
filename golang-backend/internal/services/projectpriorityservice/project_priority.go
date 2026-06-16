package projectpriorityservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectPriorityService interface {
	Add(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error)
	Update(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error)
	Delete(ctx context.Context, projectPriorityID string) error
	Get(ctx context.Context, projectPriorityID string) (domain.ProjectPriority, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error)
}

type projectPriorityService struct {
	db                        database.Database
	authorizationService      authorizationservice.AuthorizationService
	projectPriorityRepository projectpriorityrepository.ProjectPriorityRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, projectPriorityRepository projectpriorityrepository.ProjectPriorityRepository) ProjectPriorityService {
	return &projectPriorityService{db: db, authorizationService: authorizationService, projectPriorityRepository: projectPriorityRepository}
}

func (service *projectPriorityService) Add(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectPriority{}, err
	}
	projectPriority.ID = utils.UUID()
	if err := service.projectPriorityRepository.Add(ctx, service.db, projectPriority); err != nil {
		return domain.ProjectPriority{}, err
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Update(ctx context.Context, projectPriority domain.ProjectPriority) (domain.ProjectPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectPriority{}, err
	}
	if err := service.projectPriorityRepository.Update(ctx, service.db, projectPriority); err != nil {
		return domain.ProjectPriority{}, err
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Delete(ctx context.Context, projectPriorityID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.projectPriorityRepository.Delete(ctx, service.db, projectPriorityID); err != nil {
		return err
	}
	return nil
}

func (service *projectPriorityService) Get(ctx context.Context, projectPriorityID string) (domain.ProjectPriority, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.ProjectPriority{}, err
	}
	projectPriority, err := service.projectPriorityRepository.Get(ctx, service.db, projectPriorityID)
	if err != nil {
		return domain.ProjectPriority{}, fmt.Errorf("[ProjectPriorityService] failed to get project priority with ID %s: %w", projectPriorityID, err)
	}
	return projectPriority, nil
}

func (service *projectPriorityService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchProjectPrioritiesFilter) ([]domain.ProjectPriority, browser.Result, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.Result{}, err
	}
	projectPriorities, pagerResult, err := service.projectPriorityRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[ProjectPriorityService] failed to search project priorities: %w", err)
	}
	return projectPriorities, pagerResult, nil
}
