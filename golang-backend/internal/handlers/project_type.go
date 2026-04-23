package handlers

import (
	"database/sql"
	"net/http"

	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/services"
	"github.com/aportela/gotask/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectTypeHandler struct {
	service *services.ProjectTypeService
}

func NewProjectTypeHandler(service *services.ProjectTypeService) *ProjectTypeHandler {
	return &ProjectTypeHandler{service: service}
}

func (h *ProjectTypeHandler) AddProjectType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectType := models.ProjectType{}
	err := h.service.AddProjectType(ctx, projectType)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, projectType)
}

func (h *ProjectTypeHandler) UpdateProjectType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectType := models.ProjectType{}
	err := h.service.UpdateProjectType(ctx, projectType)
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
	utils.ToJSONResponse(w, http.StatusOK, projectType)
}

func (h *ProjectTypeHandler) DeleteProjectType(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectTypeId := chi.URLParam(r, "id")
	err := h.service.DeleteProjectType(ctx, projectTypeId)
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

func (h *ProjectTypeHandler) SearchProjectTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	projectTypes, err := h.service.SearchProjectTypes(ctx)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, projectTypes)
}
