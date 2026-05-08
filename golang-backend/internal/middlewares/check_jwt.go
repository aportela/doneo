package middlewares

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aportela/doneo/internal/jwt"
)

type contextKey string

const userIDKey contextKey = "userID"

func CheckJWT(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(errorResponse{UserMessage: "Authorization header missing", APIError: true}); err != nil {
					http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(errorResponse{UserMessage: "Invalid Authorization header format", APIError: true}); err != nil {
					http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			userID, err := jwt.VerifyToken(parts[1], secretKey)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				if err := json.NewEncoder(w).Encode(errorResponse{UserMessage: "Invalid JWT", ErrorDetails: err.Error(), APIError: true}); err != nil {
					http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
				}
				return
			}
			ctx := context.WithValue(r.Context(), userIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}
