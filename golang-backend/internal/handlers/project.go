package handlers

import (
	"database/sql"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
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
	project := domain.Project{}
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
	project := domain.Project{}
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

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	err := h.service.DeleteProject(ctx, projectId)
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
	utils.ToJSONResponse(w, http.StatusOK, nil)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	project, err := h.service.GetProject(ctx, projectId)
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
