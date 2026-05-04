package userhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	handlers.ToHandlerJSONResponse(w, userToAddResponse(user), nil, http.StatusCreated)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := updateRequestToUser(request)
	user.ID = chi.URLParam(r, "user_id")
	err := h.service.Update(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to update user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, userToUpdateResponse(user), nil)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "user_id")
	err := h.service.Delete(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to delete user with ID %s: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *UserHandler) Purge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "user_id")
	err := h.service.Purge(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to purge user: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "user_id")
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
	handlers.ToHandlerJSONResponse(w, userToGetResponse(user), nil)
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.service.Search(r.Context())
	handlers.ToHandlerJSONResponse(w, userArrayToSearchResponse(users), err)
}
