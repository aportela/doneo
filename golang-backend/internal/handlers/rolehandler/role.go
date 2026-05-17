package rolehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	role roleservice.RoleService
}

func NewRoleHandler(db database.Database) *RoleHandler {
	roleRepository := rolerepository.NewRoleRepository(db)
	roleService := roleservice.NewRoleService(roleRepository)
	return &RoleHandler{role: roleService}
}

func (h *RoleHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := addRequestToRole(request)
	role.ID = utils.UUID()
	err := h.role.Add(r.Context(), role)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to add role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, roleToResponse(role), nil, http.StatusCreated)
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := updateRequestToRole(request)
	role.ID = chi.URLParam(r, "id")
	err := h.role.Update(r.Context(), role)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to update role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, roleToResponse(role), nil)
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleId := chi.URLParam(r, "id")
	err := h.role.Delete(r.Context(), roleId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to delete role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleId := chi.URLParam(r, "id")
	user, err := h.role.Get(r.Context(), roleId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get non existent role: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get role: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, roleToResponse(user), nil)
}

func (h *RoleHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchRolesFilter{
		Name:                       nil,
		RequiredPermissionsBitmask: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
		if request.Filter.Permissions != nil {
			requiredPermissionsBitmask := domain.PermissionsBitmask(0)
			forbiddenPermissionsBitmask := domain.PermissionsBitmask(0)
			if request.Filter.Permissions.AllowCreate != nil {
				if *request.Filter.Permissions.AllowCreate {
					requiredPermissionsBitmask.AddPermission(domain.PermissionCreate)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionCreate)
				}
			}
			if request.Filter.Permissions.AllowUpdate != nil {
				if *request.Filter.Permissions.AllowUpdate {
					requiredPermissionsBitmask.AddPermission(domain.PermissionUpdate)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionUpdate)
				}
			}
			if request.Filter.Permissions.AllowDelete != nil {
				if *request.Filter.Permissions.AllowDelete {
					requiredPermissionsBitmask.AddPermission(domain.PermissionDelete)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionDelete)
				}
			}
			if request.Filter.Permissions.AllowView != nil {
				if *request.Filter.Permissions.AllowView {
					requiredPermissionsBitmask.AddPermission(domain.PermissionView)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionView)
				}
			}
			if request.Filter.Permissions.AllowList != nil {
				if *request.Filter.Permissions.AllowList {
					requiredPermissionsBitmask.AddPermission(domain.PermissionList)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionList)
				}
			}
			if request.Filter.Permissions.AllowExecute != nil {
				if *request.Filter.Permissions.AllowExecute {
					requiredPermissionsBitmask.AddPermission(domain.PermissionExecute)
				} else {
					forbiddenPermissionsBitmask.AddPermission(domain.PermissionExecute)
				}
			}
			filter.RequiredPermissionsBitmask = &requiredPermissionsBitmask
			filter.ForbiddenPermissionsBitmask = &forbiddenPermissionsBitmask
		}
	}

	roles, pagerResult, err := h.role.Search(r.Context(),
		browser.Params{
			CurrentPage: request.Pager.CurrentPage,
			ResultsPage: request.Pager.ResultsPage,
		},
		browser.Order{
			Field: request.Order.Field,
			Sort:  string(request.Order.Sort),
		},
		filter,
	)
	if pagerResult.TotalResults > 0 {
	}
	handlers.ToHandlerJSONResponse(w, ToSearchResponse(roles, pagerResult), err)
}
