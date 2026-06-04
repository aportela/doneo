package authhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/jwt"
	"github.com/aportela/doneo/internal/services/authservice"
	"github.com/aportela/doneo/internal/utils"
)

type AuthHandler struct {
	service                    authservice.AuthService
	secretKey                  string
	accessTokenExpirationHours int
	refreshTokenExpirationDays int
}

func NewHandler(service authservice.AuthService, secretKey string, accessTokenExpirationHours int, refreshTokenExpirationDays int) *AuthHandler {
	return &AuthHandler{service: service, secretKey: secretKey, accessTokenExpirationHours: accessTokenExpirationHours, refreshTokenExpirationDays: refreshTokenExpirationDays}
}

func (handler *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] invalid request payload: %w", err))
		return
	}
	user, err := handler.service.SignIn(r.Context(), signinRequestToUserDomain(request))
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to signin with email %s: %w", request.Email, err))
		return
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey)
	//accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(10)*time.Second), h.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate access token: %w", err))
		return
	}
	refreshToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.refreshTokenExpirationDays)*24*time.Hour), handler.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate refresh token: %w", err))
		return
	}
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken.Token,
		Path:     "/api/auth/renew-access-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshToken.ExpiresAt,
	}
	http.SetCookie(w, &refreshTokenCookie)
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken.Token,
		Path:     "/api/attachments/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshToken.ExpiresAt,
	}
	http.SetCookie(w, &accessTokenCookie)
	utils.ToJSONResponse(w, http.StatusOK,
		SignInResponse{
			AccessToken:  TokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
			RefreshToken: TokenResponse{Token: refreshToken.Token, ExpiresAt: refreshToken.ExpiresAt.UnixMilli()},
			User: userResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Permissions: userPermissions{
					IsSuperUser: user.PermissionsBitmask.HasPermission(domain.UserPermissionAdmin),
				},
			},
		},
	)
}

func (handler *AuthHandler) SignOut(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie := http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth/renew-access-token",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &refreshTokenCookie)
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/api/attachments/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &accessTokenCookie)
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *AuthHandler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	// TODO: refresh token in request ?
	if err != nil {
		// TODO: return APIError JSON HERE !
		http.Error(w, "No refresh token cookie found", http.StatusUnauthorized)
		return
	}
	var refreshToken jwt.Token
	refreshToken.Token = cookie.Value
	userId, err := jwt.VerifyToken(refreshToken.Token, handler.secretKey)
	if err != nil {
		// TODO: return APIError JSON HERE !
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	user, err := handler.service.GetUserInfo(r.Context(), userId)
	if err != nil {
		// TODO: return APIError JSON HERE !
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to get user info: %w", err))
		return
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey)
	//accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(10)*time.Second), h.secretKey)
	if err != nil {
		// TODO: return APIError JSON HERE !
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AuthHandler] failed to generate access token: %w", err))
		return
	}
	accessTokenCookie := http.Cookie{
		Name:     "access_token",
		Value:    accessToken.Token,
		Path:     "/api/attachments/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshToken.ExpiresAt,
	}
	http.SetCookie(w, &accessTokenCookie)
	utils.ToJSONResponse(w, http.StatusOK,
		RenewAccessTokenResponse{
			User: userResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Permissions: userPermissions{
					IsSuperUser: user.PermissionsBitmask.HasPermission(domain.UserPermissionAdmin),
				},
			},
			AccessToken: TokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
		},
	)
}
