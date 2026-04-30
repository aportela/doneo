package projecthandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	service projectservice.ProjectService
}

func NewProjectHandler(db database.Database) *ProjectHandler {
	projectRepository := projectrepository.NewProjectRepository(db)
	projectService := projectservice.NewProjectService(projectRepository)
	return &ProjectHandler{service: projectService}
}

func (h *ProjectHandler) AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := mapAddProjectRequestToProjectDomain(request)
	err := h.service.AddProject(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to add project with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectDomainToAddProjectResponse(project), nil, http.StatusCreated)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := mapUpdateProjectRequestToProjectDomain(request)
	err := h.service.UpdateProject(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to update project with ID %s: %w", project.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectDomainToUpdateProjectResponse(project), nil)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	err := h.service.DeleteProject(r.Context(), projectId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectService] failed to delete project with ID %s: %w", projectId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	project, err := h.service.GetProject(r.Context(), projectId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectService] not found project with ID %s: %w", projectId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectService] failed to get project with ID %s: %w", projectId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, mapProjectDomainToGetProjectResponse(project), nil)
}

func (h *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projects, err := h.service.SearchProjects(r.Context())
	handlers.ToHandlerJSONResponse(w, mapProjectArrayDomainToSearchProjectsResponse(projects), err)
}
