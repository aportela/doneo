package projectpriorityhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectPriorityHandler struct {
	service projectpriorityservice.ProjectPriorityService
}

func NewHandler(service projectpriorityservice.ProjectPriorityService) *ProjectPriorityHandler {
	return &ProjectPriorityHandler{service: service}
}

func (handler *ProjectPriorityHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] invalid request payload: %w", err))
		return
	}
	projectPriority := addRequestToDomain(request)
	projectPriority.ID = utils.UUID()
	projectPriority, err := handler.service.Add(r.Context(), projectPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to add project priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectPriority), nil, http.StatusCreated)
}

func (handler *ProjectPriorityHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] invalid request payload: %w", err))
		return
	}
	projectPriority := updateRequestToDomain(request)
	projectPriority.ID = chi.URLParam(r, "id")
	projectPriority, err := handler.service.Update(r.Context(), projectPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to update project priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectPriority), nil)
}

func (handler *ProjectPriorityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectPriorityId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), projectPriorityId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to delete project priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *ProjectPriorityHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectPriorityId := chi.URLParam(r, "id")
	projectPriority, err := handler.service.Get(r.Context(), projectPriorityId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to get non existent project priority: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] failed to get projectPriority: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectPriority), nil)
}

func (handler *ProjectPriorityHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPriorityHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectPrioritiesFilter{
		Name: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
	}
	projectPriorities, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectPriorities, pagerResult), err)
}
