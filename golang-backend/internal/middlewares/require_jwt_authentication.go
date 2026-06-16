package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/jwt"
)

func RequireJWTAuthentication(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Authorization header missing",
					"")
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Invalid Authorization header format",
					"")
				return
			}
			userID, err := jwt.VerifyToken(parts[1], secretKey)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized,
					"REQUIRE_JWT_AUTH_MIDDLEWARE_ERROR",
					"Invalid JWT",
					err.Error())
				return
			}
			//ctx := context.WithValue(r.Context(), userIDKey, userID)
			ctx := context.WithValue(r.Context(), contextUserKey, ContextUser{UserBase: domain.UserBase{ID: userID}, IsSystem: false})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
