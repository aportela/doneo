package roleservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/utils"
)

type RoleService interface {
	Add(ctx context.Context, role domain.Role) (domain.Role, error)
	Update(ctx context.Context, role domain.Role) (domain.Role, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Role, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error)
}

type roleService struct {
	database   database.Database
	repository rolerepository.RoleRepository
}

func NewService(db database.Database, role rolerepository.RoleRepository) RoleService {
	return &roleService{database: db, repository: role}
}

func (service *roleService) Add(ctx context.Context, role domain.Role) (domain.Role, error) {
	role.ID = utils.UUID()
	if err := service.repository.Add(ctx, role); err != nil {
		return domain.Role{}, fmt.Errorf("[RoleService] failed to add role with ID %s\n%w", role.ID, err)
	}
	return role, nil
}

func (service *roleService) Update(ctx context.Context, role domain.Role) (domain.Role, error) {
	if err := service.repository.Update(ctx, role); err != nil {
		return domain.Role{}, fmt.Errorf("[RoleService] failed to update role with ID %s: %w", role.ID, err)
	}
	return role, nil
}

func (service *roleService) Delete(ctx context.Context, id string) error {
	if err := service.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[RoleService] failed to delete role with ID %s: %w", id, err)
	}
	return nil
}

func (service *roleService) Get(ctx context.Context, id string) (domain.Role, error) {
	role, err := service.repository.Get(ctx, id)
	if err != nil {
		return domain.Role{}, fmt.Errorf("[RoleService] failed to get role with ID %s: %w", id, err)
	}
	return role, nil
}

func (service *roleService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error) {
	roles, pagerResult, err := service.repository.Search(ctx, pager, order, filter)
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[RoleService] failed to search roles: %w", err)
	}
	return roles, pagerResult, nil
}
