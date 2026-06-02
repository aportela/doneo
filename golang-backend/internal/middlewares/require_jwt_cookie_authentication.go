package middlewares

import (
	"context"
	"net/http"

	"github.com/aportela/doneo/internal/jwt"
)

func RequireJWTCookieAuthentication(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Access token cookie missing",
					"")
				return
			}
			userID, err := jwt.VerifyToken(cookie.Value, secretKey)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Invalid JWT",
					err.Error())
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
