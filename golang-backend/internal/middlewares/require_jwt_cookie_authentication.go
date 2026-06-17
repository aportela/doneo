package middlewares

import (
	"context"
	"net/http"

	"github.com/aportela/doneo/internal/jwt"
)

func RequireJWTCookieAuthentication(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cookie, err := r.Cookie("access_token"); err != nil {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Access token cookie missing",
					"")
				return
			} else {
				if jwtUser, err := jwt.VerifyToken(cookie.Value, secretKey); err != nil {
					writeJSONError(w, http.StatusUnauthorized,
						"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
						"Invalid JWT",
						err.Error())
					return
				} else {
					ctx := context.WithValue(r.Context(), contextUserKey, ContextUser{UserBase: jwtUser, SkipAuthorization: false})
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		})
	}
}
