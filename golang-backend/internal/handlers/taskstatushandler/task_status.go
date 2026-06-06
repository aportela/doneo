package taskstatushandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TaskStatusHandler struct {
	service taskstatusservice.TaskStatusService
}

func NewHandler(service taskstatusservice.TaskStatusService) *TaskStatusHandler {
	return &TaskStatusHandler{service: service}
}

func (handler *TaskStatusHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] invalid request payload: %w", err))
		return
	}
	taskStatus := addRequestToDomain(request)
	taskStatus.ID = utils.UUID()
	taskStatus, err := handler.service.Add(r.Context(), taskStatus)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to add project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil, http.StatusCreated)
}

func (handler *TaskStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] invalid request payload: %w", err))
		return
	}
	taskStatus := updateRequestToDomain(request)
	taskStatus.ID = chi.URLParam(r, "id")
	taskStatus, err := handler.service.Update(r.Context(), taskStatus)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to update project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil)
}

func (handler *TaskStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskStatusId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), taskStatusId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to delete project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *TaskStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskStatusId := chi.URLParam(r, "id")
	taskStatus, err := handler.service.Get(r.Context(), taskStatusId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to get non existent project status: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to get taskStatus: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil)
}

func (handler *TaskStatusHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchTaskStatusesFilter{
		Name: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
	}
	taskStatuses, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(taskStatuses, pagerResult), err)
}
