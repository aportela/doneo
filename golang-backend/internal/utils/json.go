package utils

import (
	"encoding/json"
	"net/http"
)

func ToJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "error encoding JSON: "+err.Error(), http.StatusInternalServerError)
	}
}
