package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
)

type HandlerErrorResponse struct {
	Message string `json:"message"`
	Debug   string `json:"debug"`
}

func ToHandlerJSONResponse(w http.ResponseWriter, data any, err error, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		var message string
		var httpResponseCode int
		switch {
		case errors.Is(err, domain.ErrNotFound):
			message = "not found"
			httpResponseCode = http.StatusNotFound
		case errors.Is(err, domain.ErrInvalidCredentials):
			message = "invalid credentials"
			httpResponseCode = http.StatusUnauthorized
		default:
			message = "uncaught error"
			httpResponseCode = http.StatusInternalServerError
		}
		w.WriteHeader(httpResponseCode)
		if err := json.NewEncoder(w).Encode(HandlerErrorResponse{Message: message, Debug: err.Error()}); err != nil {
			http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		code := http.StatusOK
		if len(statusCode) > 0 {
			code = statusCode[0]
		}
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
