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
	"github.com/aportela/doneo/internal/utils"
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

func (h *ProjectTypeHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	projectType := addRequestToDomain(request)
	projectType.ID = utils.UUID()
	err := h.service.Add(r.Context(), projectType)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to add project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(projectType), nil, http.StatusCreated)
}

func (h *ProjectTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	projectType := updateRequestToDomain(request)
	projectType.ID = chi.URLParam(r, "id")
	err := h.service.Update(r.Context(), projectType)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to update project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(projectType), nil)
}

func (h *ProjectTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypeId := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), projectTypeId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] failed to delete project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectTypeHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypeId := chi.URLParam(r, "id")
	projectType, err := h.service.Get(r.Context(), projectTypeId)
	if err != nil {
		if err == domain.ErrNotFound {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] failed to get non existent project type: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeService] failed to get project type: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(projectType), nil)
}

func (h *ProjectTypeHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypes, err := h.service.Search(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectTypes), err)
}
