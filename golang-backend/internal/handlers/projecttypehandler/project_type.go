package projecttypehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/go-chi/chi/v5"
)

type ProjectTypeHandler struct {
	service projecttypeservice.ProjectTypeService
}

func NewProjectTypeHandler(db database.Database) *ProjectTypeHandler {
	projectTypeRepository := projecttyperepository.NewProjectTypeRepository(db)
	projectTypeService := projecttypeservice.NewProjectTypeService(projectTypeRepository)
	return &ProjectTypeHandler{service: projectTypeService}
}

func (h *ProjectTypeHandler) AddProjectType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addProjectTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	user := mapAddProjectTypeRequestToProjectTypeDomain(request)
	err := h.service.AddProjectType(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to add project type with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectTypeDomainToAddProjectTypeResponse(user), nil, http.StatusCreated)
}

func (h *ProjectTypeHandler) UpdateProjectType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateProjectTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	user := mapUpdateProjectTypeRequestToProjectTypeDomain(request)
	err := h.service.UpdateProjectType(r.Context(), user)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to update project type with ID %s: %w", user.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectTypeDomainToUpdateProjectTypeResponse(user), nil)
}

func (h *ProjectTypeHandler) DeleteProjectType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "project_type_id")
	err := h.service.DeleteProjectType(r.Context(), userId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] failed to delete project type with ID %s: %w", userId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectTypeHandler) GetProjectType(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "project_type_id")
	user, err := h.service.GetProjectType(r.Context(), userId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] not found project type with ID %s: %w", userId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] failed to get project type with ID %s: %w", userId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, mapProjectTypeDomainToGetProjectTypeResponse(user), nil)
}

func (h *ProjectTypeHandler) SearchProjectTypes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypes, err := h.service.SearchProjectTypes(r.Context(), domain.ProjectTypeFilter{WorkspaceId: chi.URLParam(r, "workspace_id")})
	handlers.ToHandlerJSONResponse(w, mapProjectTypeArrayDomainToSearchProjectTypesResponse(projectTypes), err)
}
