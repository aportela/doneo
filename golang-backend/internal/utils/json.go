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
func ToJSONErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	if err2 := json.NewEncoder(w).Encode(map[string]string{"debugErrorMessage": err.Error()}); err2 != nil {
		http.Error(w, "error encoding JSON: "+err2.Error(), http.StatusInternalServerError)
	}
}
