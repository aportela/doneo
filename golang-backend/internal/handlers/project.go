package handlers

import (
	"database/sql"
	"net/http"

	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/services"
	"github.com/aportela/gotask/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	project := models.Project{}

	err := h.service.AddProject(ctx, project)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, project)

}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	project := models.Project{}

	err := h.service.UpdateProject(ctx, project)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ToJSONResponse(w, http.StatusNotFound, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		} else {
			utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, project)

}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	project, err := h.service.GetProject(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ToJSONResponse(w, http.StatusNotFound, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		} else {
			utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, project)

}

func (h *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	projects, err := h.service.SearchProjects(ctx)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}

	utils.ToJSONResponse(w, http.StatusOK, projects)

}
