package tasktimetrackinghandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/tasktimetrackingservice"
	"github.com/go-chi/chi/v5"
)

type TaskTimeTrackingHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type taskTimeTrackingHandler struct {
	service tasktimetrackingservice.TaskTimeTrackingService
}

func NewHandler(service tasktimetrackingservice.TaskTimeTrackingService) TaskTimeTrackingHandler {
	return &taskTimeTrackingHandler{service: service}
}

func (handler *taskTimeTrackingHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskTimeTrackingHandler] invalid request payload: %w", err))
		return
	}
	taskTimeTracking := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	if taskTimeTracking, err := handler.service.Add(r.Context(), projectID, taskID, taskTimeTracking); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskTimeTrackingHandler] failed to add task time tracking: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeTracking), nil, http.StatusCreated)
	}
}

func (handler *taskTimeTrackingHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskTimeTrackingHandler] invalid request payload: %w", err))
		return
	}
	taskTimeTracking := updateRequestToDomain(request)
	taskTimeTracking.ID = chi.URLParam(r, "task_time_tracking_id")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	if taskTimeTracking, err := handler.service.Update(r.Context(), projectID, taskID, taskTimeTracking); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskTimeTrackingHandler] failed to update task time tracking: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeTracking), nil)
	}
}

func (handler *taskTimeTrackingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	taskTimeTrackingID := chi.URLParam(r, "task_time_tracking_id")
	if err := handler.service.Delete(r.Context(), projectID, taskID, taskTimeTrackingID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TaskTimeTrackingHandler] failed to delete task time tracking: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *taskTimeTrackingHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	notes, err := handler.service.GetTaskTimeTrackings(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}
