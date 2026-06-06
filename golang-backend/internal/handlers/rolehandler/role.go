package rolehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/go-chi/chi/v5"
)

type RoleHandler struct {
	service roleservice.RoleService
}

func NewHandler(service roleservice.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (handler *RoleHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := addRequestToDomain(request)
	role, err := handler.service.Add(r.Context(), role)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to add role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(role), nil, http.StatusCreated)
}

func (handler *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := updateRequestToDomain(request)
	role.ID = chi.URLParam(r, "id")
	role, err := handler.service.Update(r.Context(), role)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to update role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(role), nil)
}

func (handler *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), roleId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to delete role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *RoleHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleId := chi.URLParam(r, "id")
	user, err := handler.service.Get(r.Context(), roleId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to get non existent role: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to get role: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(user), nil)
}

func (handler *RoleHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roles, _, err := handler.service.Search(r.Context(),
		browser.Params{
			CurrentPage: 1,
			ResultsPage: 0,
		},
		browser.Order{
			Field: "name",
			Sort:  "ASC",
		},
		domain.SearchRolesFilter{},
	)
	handlers.ToHandlerJSONResponse(w, toSearchBaseResponse(roles), err)
}

func (handler *RoleHandler) Search(w http.ResponseWriter, r *http.Request) {
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
	}

	roles, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(roles, pagerResult), err)
}
