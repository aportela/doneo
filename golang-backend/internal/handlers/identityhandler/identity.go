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

type IdentityHandler interface {
	SignIn(w http.ResponseWriter, r *http.Request)
	SignOut(w http.ResponseWriter, r *http.Request)
	RenewAccessToken(w http.ResponseWriter, r *http.Request)
}

type identityHandler struct {
	identityService            identityservice.IdentityService
	secretKey                  string
	accessTokenExpirationHours int
	refreshTokenExpirationDays int
}

func NewHandler(identityService identityservice.IdentityService, secretKey string, accessTokenExpirationHours int, refreshTokenExpirationDays int) IdentityHandler {
	return &identityHandler{identityService: identityService, secretKey: secretKey, accessTokenExpirationHours: accessTokenExpirationHours, refreshTokenExpirationDays: refreshTokenExpirationDays}
}

func (handler *identityHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var request signInRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] invalid request payload: %w", err))
		return
	}
	if user, err := handler.identityService.SignIn(r.Context(), request.Email, request.Password); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to signin with email %s: %w", request.Email, err))
		return
	} else {
		if accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate access token: %w", err))
			return
		} else {
			if refreshToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.refreshTokenExpirationDays)*24*time.Hour), handler.secretKey); err != nil {
				handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate refresh token: %w", err))
				return
			} else {
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
					Path:     "/api/wc/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteStrictMode,
					Expires:  refreshToken.ExpiresAt,
				}
				http.SetCookie(w, &accessTokenCookie)
				utils.ToJSONResponse(w, http.StatusOK,
					signInResponse{
						AccessToken:  tokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
						RefreshToken: tokenResponse{Token: refreshToken.Token, ExpiresAt: refreshToken.ExpiresAt.UnixMilli()},
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
		}
	}
}

func (handler *identityHandler) SignOut(w http.ResponseWriter, r *http.Request) {
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
		Path:     "/api/wc/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-time.Hour),
	}
	http.SetCookie(w, &accessTokenCookie)
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *identityHandler) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("refresh_token"); err != nil {
		// TODO: return APIError JSON HERE !
		http.Error(w, "No refresh token cookie found", http.StatusUnauthorized)
		return
	} else {
		var refreshToken jwt.Token
		refreshToken.Token = cookie.Value
		if jwtUser, err := jwt.VerifyToken(refreshToken.Token, handler.secretKey); err != nil {
			// TODO: return APIError JSON HERE !
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		} else {
			ctx := middlewares.SetContextUser(r.Context(), middlewares.ContextUser{UserBase: jwtUser, SkipAuthorization: false})
			if user, err := handler.identityService.GetCurrentUserInfo(ctx); err != nil {
				// TODO: return APIError JSON HERE !
				handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to get user info: %w", err))
				return
			} else {
				if accessToken, err := jwt.GenerateToken(user, time.Now().Add(time.Duration(handler.accessTokenExpirationHours)*time.Hour), handler.secretKey); err != nil {
					// TODO: return APIError JSON HERE !
					handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[IdentityHandler] failed to generate access token: %w", err))
					return
				} else {
					accessTokenCookie := http.Cookie{
						Name:     "access_token",
						Value:    accessToken.Token,
						Path:     "/api/wc/",
						HttpOnly: true,
						Secure:   true,
						SameSite: http.SameSiteStrictMode,
						Expires:  refreshToken.ExpiresAt,
					}
					http.SetCookie(w, &accessTokenCookie)
					utils.ToJSONResponse(w, http.StatusOK,
						renewAccessTokenResponse{
							User: userResponse{
								ID:    user.ID,
								Name:  user.Name,
								Email: user.Email,
								Permissions: userPermissions{
									IsSuperUser: user.PermissionsBitmask.HasFlag(domain.UserPermissionAdmin),
								},
							},
							AccessToken: tokenResponse{Token: accessToken.Token, ExpiresAt: accessToken.ExpiresAt.UnixMilli()},
						},
					)
				}
			}
		}
	}
}
