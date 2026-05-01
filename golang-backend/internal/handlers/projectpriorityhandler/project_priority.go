package projectpriorityhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/go-chi/chi/v5"
)

type ProjectPriorityHandler struct {
	service projectpriorityservice.ProjectPriorityService
}

func NewProjectPriorityHandler(db database.Database) *ProjectPriorityHandler {
	projectPriorityRepository := projectpriorityrepository.NewProjectPriorityRepository(db)
	projectPriorityService := projectpriorityservice.NewProjectPriorityService(projectPriorityRepository)
	return &ProjectPriorityHandler{service: projectPriorityService}
}

func (h *ProjectPriorityHandler) AddProjectPriority(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addProjectPriorityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] invalid request payload: %w", err))
		return
	}
	projectPriority := mapAddProjectPriorityRequestToProjectPriorityDomain(request)
	err := h.service.AddProjectPriority(r.Context(), projectPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to add project priority with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectPriorityDomainToAddProjectPriorityResponse(projectPriority), nil, http.StatusCreated)
}

func (h *ProjectPriorityHandler) UpdateProjectPriority(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateProjectPriorityRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] invalid request payload: %w", err))
		return
	}
	projectPriority := mapUpdateProjectPriorityRequestToProjectPriorityDomain(request)
	err := h.service.UpdateProjectPriority(r.Context(), projectPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to update project priority with ID %s: %w", projectPriority.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, mapProjectPriorityDomainToUpdateProjectPriorityResponse(projectPriority), nil)
}

func (h *ProjectPriorityHandler) DeleteProjectPriority(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectPriorityId := chi.URLParam(r, "id")
	err := h.service.DeleteProjectPriority(r.Context(), projectPriorityId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityService] failed to delete project priority with ID %s: %w", projectPriorityId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectPriorityHandler) GetProjectPriority(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectPriorityId := chi.URLParam(r, "id")
	projectPriority, err := h.service.GetProjectPriority(r.Context(), projectPriorityId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityService] not found projectPriority with ID %s: %w", projectPriorityId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityService] failed to get projectPriority with ID %s: %w", projectPriorityId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, mapProjectPriorityDomainToGetProjectPriorityResponse(projectPriority), nil)
}

func (h *ProjectPriorityHandler) SearchProjectPriorities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectPriorities, err := h.service.SearchProjectPriorities(r.Context())
	handlers.ToHandlerJSONResponse(w, mapProjectPriorityArrayDomainToSearchProjectPrioritysResponse(projectPriorities), err)
}
