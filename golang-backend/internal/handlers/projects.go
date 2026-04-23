package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aportela/gotask/internal/services"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	project, err := h.service.Get(id)
	if err != nil {
		http.Error(w, "error"+err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(project)

}
