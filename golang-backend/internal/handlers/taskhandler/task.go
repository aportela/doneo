package taskhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/taskservice"
	"github.com/go-chi/chi/v5"
)

type TaskHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
	service taskservice.TaskService
}

func NewHandler(service taskservice.TaskService) TaskHandler {
	return &taskHandler{service: service}
}

func (handler *taskHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	task := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	task, err := handler.service.Add(r.Context(), projectID, task)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to add task with ID %s: %w", request.ID, err))
		return
	}
	task, err = handler.service.Get(r.Context(), projectID, task.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get new task with ID %s: %w", task.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(task), nil, http.StatusCreated)
}

func (handler *taskHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	task := updateRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	task.ID = chi.URLParam(r, "task_id")
	task, err := handler.service.Update(r.Context(), projectID, task)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to update task with ID %s: %w", task.ID, err))
		return
	}
	task, err = handler.service.Get(r.Context(), projectID, task.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get updated task with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(task), nil)
}

func (handler *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	err := handler.service.Delete(r.Context(), projectID, taskID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to delete project with ID %s: %w", projectID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *taskHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	project, err := handler.service.Get(r.Context(), projectID, taskID)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] not found project with ID %s: %w", taskID, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get project with ID %s: %w", taskID, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (handler *taskHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	projectID := chi.URLParam(r, "project_id")
	filter := domain.SearchTaskFilter{}
	filter.ProjectID = &projectID

	if request.Filter != nil {
		if request.Filter.Summary != nil {
			filter.Summary = request.Filter.Summary
		}
		if request.Filter.PriorityID != nil {
			filter.PriorityID = request.Filter.PriorityID
		}
		if request.Filter.StatusID != nil {
			filter.StatusID = request.Filter.StatusID
		}
		if request.Filter.CreatedAt != nil {
			filter.CreatedAt = &domain.TimestampFilter{From: nil, To: nil}
			if request.Filter.CreatedAt.From != nil {
				filter.CreatedAt.From = request.Filter.CreatedAt.From
			}
			if request.Filter.CreatedAt.To != nil {
				filter.CreatedAt.To = request.Filter.CreatedAt.To
			}
		}
		if request.Filter.CreatedByUserID != nil {
			filter.CreatedByUserID = request.Filter.CreatedByUserID
		}
	}
	projects, pagerResult, err := handler.service.Search(r.Context(),
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
