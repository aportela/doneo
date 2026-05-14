package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/userservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service userservice.UserService
}

func NewUserHandler(db database.Database) *UserHandler {
	userRepository := userrepository.NewUserRepository(db)
	userService := userservice.NewUserService(userRepository)
	return &UserHandler{service: userService}
}

func (h *UserHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := addRequestToUser(request)
	user.ID = utils.UUID()
	err := h.service.Add(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to add user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, userToResponse(user), nil, http.StatusCreated)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := updateRequestToUser(request)
	user.ID = chi.URLParam(r, "id")
	err := h.service.Update(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to update user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, userToResponse(user), nil)
}

func (h *UserHandler) Patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request patchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	userId := chi.URLParam(r, "id")

	user, err := h.service.Get(r.Context(), userId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get non existent user: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get user: %w", err))
			return
		}
	}
	if request.DeletedAt == nil {
		user.DeletedAt = nil
	}
	err = h.service.Patch(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to patch user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, userToResponse(user), nil)

}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	// TODO: deny delete current session user ?
	err := h.service.Delete(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to delete user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *UserHandler) Purge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	err := h.service.Purge(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to purge user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	user, err := h.service.Get(r.Context(), userId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get non existent user: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get user: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, userToResponse(user), nil)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchUsersFilter{
		Name:              nil,
		Email:             nil,
		AdministratorFlag: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
		if request.Filter.Email != nil {
			filter.Email = request.Filter.Email
		}
		if request.Filter.AdministratorFlag != nil {
			filter.AdministratorFlag = request.Filter.AdministratorFlag
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
	users, pagerResult, err := h.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, ToSearchResponse(users, pagerResult), err)
}
