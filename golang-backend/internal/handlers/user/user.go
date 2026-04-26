package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(db database.Database) *UserHandler {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	return &UserHandler{service: userService}
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := ToUser(userRequest)
	err := h.service.AddUser(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to add user with ID %s: %w", userRequest.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, ToUserResponse(user), nil, http.StatusCreated)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userRequest UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	user := ToUser(userRequest)
	err := h.service.UpdateUser(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] failed to update user with ID %s: %w", user.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, ToUserResponse(user), nil)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	err := h.service.DeleteUser(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to delete user with ID %s: %w", userId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	user, err := h.service.GetUser(r.Context(), userId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] not found user with ID %s: %w", userId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserService] failed to get user with ID %s: %w", userId, err))
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, ToUserResponse(user))
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := h.service.SearchUsers(r.Context())
	handlers.ToHandlerJSONResponse(w, ToSearchUserResponse(users), err)
}
