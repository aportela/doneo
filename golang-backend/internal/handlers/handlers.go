package handlers

import (
	"encoding/json"
	"net/http"
)

type HandlerErrorResponse struct {
	Message  string `json:"message"`
	Debug    string `json:"debug"`
	APIError bool   `json:"APIError"`
	Details  any    `json:"details"`
}

func ToHandlerJSONResponse(w http.ResponseWriter, data any, err error, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {

		status, message, details := mapError(err)

		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(HandlerErrorResponse{Message: message, Debug: err.Error(), APIError: true, Details: details}); err != nil {
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

type TimestampFilter struct {
	From *uint64 `json:"from"`
	To   *uint64 `json:"to"`
}
