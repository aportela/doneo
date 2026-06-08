package projecttaskhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projecttaskservice"
	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	service projecttaskservice.TaskService
}

func NewHandler(service projecttaskservice.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (handler *TaskHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	task := addRequestToDomain(request)
	projectId := chi.URLParam(r, "id")
	task, err := handler.service.Add(r.Context(), projectId, task)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to add task with ID %s: %w", request.ID, err))
		return
	}
	task, err = handler.service.Get(r.Context(), task.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get new task with ID %s: %w", task.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(task), nil, http.StatusCreated)
}

func (handler *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	task := updateRequestToDomain(request)
	task.ID = chi.URLParam(r, "task_id")
	task, err := handler.service.Update(r.Context(), task)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to update task with ID %s: %w", task.ID, err))
		return
	}
	task, err = handler.service.Get(r.Context(), task.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get updated task with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(task), nil)
}

func (handler *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), projectId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to delete project with ID %s: %w", projectId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskId := chi.URLParam(r, "task_id")
	project, err := handler.service.Get(r.Context(), taskId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] not found project with ID %s: %w", taskId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] failed to get project with ID %s: %w", taskId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (handler *TaskHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskHandler] invalid request payload: %w", err))
		return
	}
	projectId := chi.URLParam(r, "id")
	filter := domain.SearchTaskFilter{}
	filter.ProjectId = &projectId

	if request.Filter != nil {
		if request.Filter.Summary != nil {
			filter.Summary = request.Filter.Summary
		}
		if request.Filter.PriorityId != nil {
			filter.PriorityId = request.Filter.PriorityId
		}
		if request.Filter.StatusId != nil {
			filter.StatusId = request.Filter.StatusId
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
		if request.Filter.CreatedByUserId != nil {
			filter.CreatedByUserId = request.Filter.CreatedByUserId
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
