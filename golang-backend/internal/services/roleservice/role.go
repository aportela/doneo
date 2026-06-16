package roleservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/utils"
)

type RoleService interface {
	Add(ctx context.Context, role domain.Role) (domain.Role, error)
	Update(ctx context.Context, role domain.Role) (domain.Role, error)
	Delete(ctx context.Context, roleID string) error
	Get(ctx context.Context, roleID string) (domain.Role, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error)
}

type roleService struct {
	db                   database.Database
	authorizationService authorizationservice.AuthorizationService
	roleRepository       rolerepository.RoleRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, roleRepository rolerepository.RoleRepository) RoleService {
	return &roleService{db: db, authorizationService: authorizationService, roleRepository: roleRepository}
}

func (service *roleService) Add(ctx context.Context, role domain.Role) (domain.Role, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.Role{}, err
	}
	role.ID = utils.UUID()
	if err := service.roleRepository.Add(ctx, service.db, role); err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (service *roleService) Update(ctx context.Context, role domain.Role) (domain.Role, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.Role{}, err
	}
	if err := service.roleRepository.Update(ctx, service.db, role); err != nil {
		return domain.Role{}, err
	}
	return role, nil
}

func (service *roleService) Delete(ctx context.Context, roleID string) error {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return err
	}
	if err := service.roleRepository.Delete(ctx, service.db, roleID); err != nil {
		return err
	}
	return nil
}

func (service *roleService) Get(ctx context.Context, roleID string) (domain.Role, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return domain.Role{}, err
	}
	role, err := service.roleRepository.Get(ctx, service.db, roleID)
	if err != nil {
		return domain.Role{}, fmt.Errorf("[RoleService] failed to get role with ID %s: %w", roleID, err)
	}
	return role, nil
}

func (service *roleService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error) {
	if _, err := service.authorizationService.RequireUserAdminPermission(ctx); err != nil {
		return nil, browser.Result{}, err
	}
	roles, pagerResult, err := service.roleRepository.Search(ctx, service.db, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[RoleService] failed to search roles: %w", err)
	}
	return roles, pagerResult, nil
}
