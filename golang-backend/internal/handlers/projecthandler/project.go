package projecthandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/utils"
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

func (h *ProjectHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := addRequestToDomain(request)
	project.ID = utils.UUID()
	project.CreatedBy = domain.UserBase{}
	project.CreatedBy.ID, _ = middlewares.GetUserIDFromContext(r.Context())
	err := h.service.Add(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to add project with ID %s: %w", request.ID, err))
		return
	}
	project, err = h.service.Get(r.Context(), project.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get new project with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil, http.StatusCreated)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := updateRequestToDomain(request)
	project.ID = chi.URLParam(r, "id")
	err := h.service.Update(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to update project with ID %s: %w", project.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), projectId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to delete project with ID %s: %w", projectId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	project, err := h.service.Get(r.Context(), projectId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] not found project with ID %s: %w", projectId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get project with ID %s: %w", projectId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (h *ProjectHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectFilter{
		Key: nil,
	}
	if request.Filter != nil {
		if request.Filter.Key != nil {
			filter.Key = request.Filter.Key
		}
	}
	projects, pagerResult, err := h.service.Search(r.Context(),
		browser.Params{
			CurrentPage: request.Pager.CurrentPage,
			ResultsPage: request.Pager.ResultsPage,
		},
		browser.Order{
			Field: request.Order.Field,
			Sort:  string(request.Order.Sort),
		},
		filter,
	)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projects, pagerResult), err)
}
