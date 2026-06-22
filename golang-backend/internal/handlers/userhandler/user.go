package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/userservice"
	"github.com/go-chi/chi/v5"
)

type UserHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Purge(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	SearchBase(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	service userservice.UserService
}

func NewHandler(service userservice.UserService) UserHandler {
	return &userHandler{service: service}
}

func (handler *userHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := addRequestToDomain(request)
	if user, err := handler.service.Add(r.Context(), user, request.Password); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to add user: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil, http.StatusCreated)
	}
}

func (handler *userHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := updateRequestToDomain(request)
	user.ID = chi.URLParam(r, "id")
	if user, err := handler.service.Update(r.Context(), user, request.Password); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to update user: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil)
	}
}

func (handler *userHandler) Patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request patchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	userID := chi.URLParam(r, "id")
	if user, err := handler.service.Get(r.Context(), userID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to get non existent user: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to get user: %w", err))
			return
		}
	} else {
		if request.DeletedAt == nil {
			user.DeletedAt = nil
		}
		if err := handler.service.Patch(r.Context(), user); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to patch user: %w", err))
			return
		}
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil)
	}
}

func (handler *userHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	if err := handler.service.Delete(r.Context(), userID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to delete user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *userHandler) Purge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	if err := handler.service.Purge(r.Context(), userID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to purge user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "id")
	if user, err := handler.service.Get(r.Context(), userID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get non existent user: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get user: %w", err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil)
	}
}

func (handler *userHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, _, err := handler.service.Search(r.Context(),
		browser.Params{
			CurrentPage: 1,
			ResultsPage: 0,
		},
		browser.Order{
			Field: "name",
			Sort:  "ASC",
		},
		domain.SearchUsersFilter{},
	)
	handlers.ToHandlerJSONResponse(w, toSearchBaseResponse(users), err)
}

func (handler *userHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchUsersFilter{
		Name:                        nil,
		Email:                       nil,
		RequiredPermissionsBitmask:  nil,
		ForbiddenPermissionsBitmask: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
		if request.Filter.Email != nil {
			filter.Email = request.Filter.Email
		}
		if request.Filter.Permissions != nil {
			requiredPermissionsBitmask := domain.Bitmask(0)
			forbiddenPermissionsBitmask := domain.Bitmask(0)
			if request.Filter.Permissions.IsSuperUser != nil {
				if *request.Filter.Permissions.IsSuperUser {
					requiredPermissionsBitmask.AddFlag(domain.UserPermissionAdmin)
					filter.RequiredPermissionsBitmask = &requiredPermissionsBitmask
				} else {
					forbiddenPermissionsBitmask.AddFlag(domain.UserPermissionAdmin)
					filter.ForbiddenPermissionsBitmask = &forbiddenPermissionsBitmask
				}
			}
		}
		if request.Filter.CreatedAt != nil {
			filter.CreatedAt = &domain.TimestampFilter{From: nil, To: nil}
			if request.Filter.CreatedAt.From != nil {
				filter.CreatedAt.From = request.Filter.CreatedAt.From
			}
			if request.Filter.CreatedAt.To != nil {
				filter.CreatedAt.To = request.Filter.CreatedAt.To
			}
		}
		if request.Filter.UpdatedAt != nil {
			filter.UpdatedAt = &domain.TimestampFilter{From: nil, To: nil}
			if request.Filter.UpdatedAt.From != nil {
				filter.UpdatedAt.From = request.Filter.UpdatedAt.From
			}
			if request.Filter.CreatedAt.To != nil {
				filter.UpdatedAt.To = request.Filter.UpdatedAt.To
			}
		}
		if request.Filter.DeletedAt != nil {
			filter.DeletedAt = &domain.TimestampFilter{From: nil, To: nil}
			if request.Filter.DeletedAt.From != nil {
				filter.DeletedAt.From = request.Filter.DeletedAt.From
			}
			if request.Filter.DeletedAt.To != nil {
				filter.DeletedAt.To = request.Filter.DeletedAt.To
			}
		}
	}
	users, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(users, pagerResult), err)
}
