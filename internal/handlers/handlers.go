package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool `json:"success"`
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := Response{
		Success: true,
	}

	json.NewEncoder(w).Encode(resp)

}
