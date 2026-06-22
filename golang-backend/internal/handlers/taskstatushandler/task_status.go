package taskstatushandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/go-chi/chi/v5"
)

type TaskStatusHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	SearchBase(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type taskStatusHandler struct {
	service taskstatusservice.TaskStatusService
}

func NewHandler(service taskstatusservice.TaskStatusService) TaskStatusHandler {
	return &taskStatusHandler{service: service}
}

func (handler *taskStatusHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] invalid request payload: %w", err))
		return
	}
	taskStatus := addRequestToDomain(request)
	if taskStatus, err := handler.service.Add(r.Context(), taskStatus); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to add project status: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil, http.StatusCreated)
	}
}

func (handler *taskStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] invalid request payload: %w", err))
		return
	}
	taskStatus := updateRequestToDomain(request)
	taskStatus.ID = chi.URLParam(r, "id")
	if taskStatus, err := handler.service.Update(r.Context(), taskStatus); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to update project status: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil)
	}
}

func (handler *taskStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskStatusID := chi.URLParam(r, "id")
	if err := handler.service.Delete(r.Context(), taskStatusID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to delete project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *taskStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskStatusID := chi.URLParam(r, "id")
	if taskStatus, err := handler.service.Get(r.Context(), taskStatusID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to get non existent project status: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskStatusHandler] failed to get taskStatus: %w", err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(taskStatus), nil)
	}
}

func (handler *taskStatusHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskStatuses, pagerResult, err := handler.service.SearchBase(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(taskStatuses, pagerResult), err)
}

func (handler *taskStatusHandler) Search(w http.ResponseWriter, r *http.Request) {
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
