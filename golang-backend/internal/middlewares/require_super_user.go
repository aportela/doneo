package middlewares

import (
	"encoding/json"
	"net/http"
)

func isSuperUser() bool {
	return true
}

func RequireSuperUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isSuperUser() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			resp := errorResponse{
				Message: "missing super user flag",
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
