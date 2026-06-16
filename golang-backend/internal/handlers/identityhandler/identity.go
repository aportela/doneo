package identityhandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/jwt"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/services/identityservice"
	"github.com/aportela/doneo/internal/utils"
)

type IdentityHandler struct {
	identityService            identityservice.IdentityService
	secretKey                  string
	accessTokenExpirationHours int
	refreshTokenExpirationDays int
}

// TODO: add func definitions

func NewHandler(identityService identityservice.IdentityService, secretKey string, accessTokenExpirationHours int, refreshTokenExpirationDays int) *IdentityHandler {
	return &IdentityHandler{identityService: identityService, secretKey: secretKey, accessTokenExpirationHours: accessTokenExpirationHours, refreshTokenExpirationDays: refreshTokenExpirationDays}
}

func (handler *IdentityHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] invalid request payload: %w", err))
		return
	}
	user, err := handler.identityService.SignIn(r.Context(), request.Email, request.Password)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to signin with email %s: %w", request.Email, err))
		return
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey)
	//accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(10)*time.Second), h.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate access token: %w", err))
		return
	}
	refreshToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.refreshTokenExpirationDays)*24*time.Hour), handler.secretKey)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate refresh token: %w", err))
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
					IsSuperUser: user.PermissionsBitmask.HasFlag(domain.UserPermissionAdmin),
				},
			},
		},
	)
}

func (handler *IdentityHandler) SignOut(w http.ResponseWriter, r *http.Request) {
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

func (handler *IdentityHandler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	// TODO: refresh token in request ?
	if err != nil {
		// TODO: return APIError JSON HERE !
		http.Error(w, "No refresh token cookie found", http.StatusUnauthorized)
		return
	}
	var refreshToken jwt.Token
	refreshToken.Token = cookie.Value
	userID_, err := jwt.VerifyToken(refreshToken.Token, handler.secretKey)
	if err != nil {
		// TODO: return APIError JSON HERE !
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//ctx := middlewares.SetUserIDIntoContext(r.Context(), userID_)
	ctx := middlewares.SetContextUser(r.Context(), middlewares.ContextUser{UserBase: domain.UserBase{ID: userID_}, IsSystem: false})
	user, err := handler.identityService.GetCurrentUserInfo(ctx)
	if err != nil {
		// TODO: return APIError JSON HERE !
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to get user info: %w", err))
		return
	}
	accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey)
	//accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(10)*time.Second), h.secretKey)
	if err != nil {
		// TODO: return APIError JSON HERE !
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate access token: %w", err))
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
					IsSuperUser: user.PermissionsBitmask.HasFlag(domain.UserPermissionAdmin),
				},
			},
			AccessToken: TokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
		},
	)
}
