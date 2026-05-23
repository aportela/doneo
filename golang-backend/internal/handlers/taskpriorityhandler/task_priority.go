package taskpriorityhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TaskPriorityHandler struct {
	service taskpriorityservice.TaskPriorityService
}

func NewTaskPriorityHandler(db database.Database) *TaskPriorityHandler {
	taskPriorityRepository := taskpriorityrepository.NewTaskPriorityRepository(db)
	taskPriorityService := taskpriorityservice.NewTaskPriorityService(taskPriorityRepository)
	return &TaskPriorityHandler{service: taskPriorityService}
}

func (h *TaskPriorityHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
		return
	}
	taskPriority := addRequestToDomain(request)
	taskPriority.ID = utils.UUID()
	err := h.service.Add(r.Context(), taskPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to add task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(taskPriority), nil, http.StatusCreated)
}

func (h *TaskPriorityHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
		return
	}
	taskPriority := updateRequestToDomain(request)
	taskPriority.ID = chi.URLParam(r, "id")
	err := h.service.Update(r.Context(), taskPriority)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to update task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(taskPriority), nil)
}

func (h *TaskPriorityHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskPriorityId := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), taskPriorityId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to delete task priority: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *TaskPriorityHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskPriorityId := chi.URLParam(r, "id")
	taskPriority, err := h.service.Get(r.Context(), taskPriorityId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to get non existent task priority: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] failed to get task priority: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(taskPriority), nil)
}

func (h *TaskPriorityHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskPriorityHandler] invalid request payload: %w", err))
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
	taskPriorities, pagerResult, err := h.service.Search(r.Context(),
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
