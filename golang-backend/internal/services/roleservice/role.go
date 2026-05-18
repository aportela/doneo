package roleservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
)

type RoleService interface {
	Add(ctx context.Context, role domain.Role) error
	Update(ctx context.Context, role domain.Role) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.Role, error)
	Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error)
}

type roleService struct {
	repository rolerepository.RoleRepository
}

func NewRoleService(role rolerepository.RoleRepository) RoleService {
	return &roleService{repository: role}
}

func (r *roleService) Add(ctx context.Context, role domain.Role) error {
	if err := r.repository.Add(ctx, rolerepository.RoleToDTO(role)); err != nil {
		return fmt.Errorf("[RoleService] failed to add role with ID %s: %w", role.ID, err)
	}
	return nil
}

func (r *roleService) Update(ctx context.Context, role domain.Role) error {
	if err := r.repository.Update(ctx, rolerepository.RoleToDTO(role)); err != nil {
		return fmt.Errorf("[RoleService] failed to update role with ID %s: %w", role.ID, err)
	}
	return nil
}

func (r *roleService) Delete(ctx context.Context, id string) error {
	if err := r.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[RoleService] failed to delete role with ID %s: %w", id, err)
	}
	return nil
}

func (r *roleService) Get(ctx context.Context, id string) (domain.Role, error) {
	role, err := r.repository.Get(ctx, id)
	if err != nil {
		return domain.Role{}, fmt.Errorf("[RoleService] failed to get role with ID %s: %w", id, err)
	}
	return rolerepository.DTOToDomain(role), nil
}

func (r *roleService) Search(ctx context.Context, pager browser.Params, order browser.Order, filter domain.SearchRolesFilter) ([]domain.Role, browser.Result, error) {
	roles, pagerResult, err := r.repository.Search(ctx, pager, order, rolerepository.SearchRolesFilterToDTO(filter))
	if err != nil {
		return nil, browser.Result{}, fmt.Errorf("[RoleService] failed to search roles: %w", err)
	}
	return rolerepository.DTOArrayToDomainArray(roles), pagerResult, nil
}
