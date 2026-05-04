package authhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/jwt"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authservice"
	"github.com/aportela/doneo/internal/utils"
)

type AuthHandler struct {
	service                    authservice.AuthService
	secretKey                  string
	accessTokenExpirationHours int
	refreshTokenExpirationDays int
}

func NewAuthHandler(db database.Database, secretKey string, accessTokenExpirationHours int, refreshTokenExpirationDays int) *AuthHandler {
	userRepository := userrepository.NewUserRepository(db)
	authService := authservice.NewAuthService(userRepository)
	return &AuthHandler{service: authService, secretKey: secretKey, accessTokenExpirationHours: accessTokenExpirationHours, refreshTokenExpirationDays: refreshTokenExpirationDays}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] invalid request payload: %w", err))
		return
	}
	user, err := h.service.SignIn(r.Context(), signinRequestToUserDomain(request))
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to signin with email %s: %w", request.Email, err))
		return
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(h.accessTokenExpirationHours)*time.Hour), h.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate access token: %w", err))
		return
	}
	refreshToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(h.refreshTokenExpirationDays)*24*time.Hour), h.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate refresh token: %w", err))
		return
	}
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Token,
		Path:     "/api/auth/renew-access-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshToken.ExpiresAt,
	}
	http.SetCookie(w, &cookie)
	utils.ToJSONResponse(w, http.StatusOK,
		SignInResponse{
			AccessToken:  TokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
			RefreshToken: TokenResponse{Token: refreshToken.Token, ExpiresAt: refreshToken.ExpiresAt.UnixMilli()},
		},
	)
}

func (h *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/renew-access-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &cookie)
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (h *AuthHandler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	// TODO: refresh token in request ?
	if err != nil {
		http.Error(w, "No refresh token cookie found", http.StatusUnauthorized)
		return
	}
	var refreshToken jwt.Token
	refreshToken.Token = cookie.Value
	userId, err := jwt.VerifyToken(refreshToken.Token, h.secretKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	user := domain.User{
		UserBase: domain.UserBase{ID: userId},
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(h.accessTokenExpirationHours)*time.Hour), h.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate access token: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK,
		RenewAccessTokenResponse{
			AccessToken: TokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
		},
	)
}
