package taskpriorityhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/go-chi/chi/v5"
)

type TaskPriorityHandler struct {
	service taskpriorityservice.TaskPriorityService
}

func NewHandler(service taskpriorityservice.TaskPriorityService) *TaskPriorityHandler {
	return &TaskPriorityHandler{service: service}
}

func (handler *TaskPriorityHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
		return
	}
	taskPriority := addRequestToDomain(request)
	taskPriority, err := handler.service.Add(r.Context(), taskPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to add task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskPriority), nil, http.StatusCreated)
}

func (handler *TaskPriorityHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
		return
	}
	taskPriority := updateRequestToDomain(request)
	taskPriority.ID = chi.URLParam(r, "id")
	taskPriority, err := handler.service.Update(r.Context(), taskPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to update task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskPriority), nil)
}

func (handler *TaskPriorityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskPriorityId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), taskPriorityId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to delete task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *TaskPriorityHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskPriorityId := chi.URLParam(r, "id")
	taskPriority, err := handler.service.Get(r.Context(), taskPriorityId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to get non existent task priority: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to get task priority: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskPriority), nil)
}

func (handler *TaskPriorityHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchTaskPrioritiesFilter{
		Name: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
	}
	taskPriorities, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(taskPriorities, pagerResult), err)
}
