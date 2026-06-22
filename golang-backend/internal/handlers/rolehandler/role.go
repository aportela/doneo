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

type RoleHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	SearchBase(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type roleHandler struct {
	service roleservice.RoleService
}

func NewHandler(service roleservice.RoleService) RoleHandler {
	return &roleHandler{service: service}
}

func (handler *roleHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := addRequestToDomain(request)
	if role, err := handler.service.Add(r.Context(), role); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to add role: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(role), nil, http.StatusCreated)
	}
}

func (handler *roleHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] invalid request payload: %w", err))
		return
	}
	role := updateRequestToDomain(request)
	role.ID = chi.URLParam(r, "id")
	if role, err := handler.service.Update(r.Context(), role); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to update role: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(role), nil)
	}
}

func (handler *roleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleID := chi.URLParam(r, "id")
	if err := handler.service.Delete(r.Context(), roleID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to delete role: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *roleHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	roleID := chi.URLParam(r, "id")
	if user, err := handler.service.Get(r.Context(), roleID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to get non existent role: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[RoleHandler] failed to get role: %w", err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(user), nil)
	}
}

func (handler *roleHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
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

func (handler *roleHandler) Search(w http.ResponseWriter, r *http.Request) {
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
