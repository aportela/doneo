package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(db database.Database, secretKey string) *AuthHandler {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository, secretKey)
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&signUpRequest); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] invalid request payload: %w", err))
		return
	}
	err := h.service.SignUp(r.Context(), ToUser(signUpRequest))
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to register user with email %s: %w", signUpRequest.Email, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil, http.StatusCreated)
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var signInRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&signInRequest); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] invalid request payload: %w", err))
		return
	}
	token, err := h.service.SignIn(r.Context(), signInRequest.Email, signInRequest.Password)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to signin with email %s: %w", signInRequest.Email, err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, SuccessSignInResponse{Token: token})

}
