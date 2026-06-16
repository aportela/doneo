package tasktimerentryhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/tasktimerentryservice"
	"github.com/go-chi/chi/v5"
)

type TaskTimerEntryHandler struct {
	service tasktimerentryservice.TaskTimerEntryService
}

func NewHandler(service tasktimerentryservice.TaskTimerEntryService) *TaskTimerEntryHandler {
	return &TaskTimerEntryHandler{service: service}
}

func (handler *TaskTimerEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	taskTimeEntry := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")

	// TODO: return taskTimeEntry with new id & createdAt
	err := handler.service.Add(r.Context(), projectID, taskID, taskTimeEntry)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add task time entry: %w", err))
		return
	}

	handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeEntry), nil, http.StatusCreated)
}

func (handler *TaskTimerEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	taskTimeEntry := updateRequestToDomain(request)
	taskTimeEntry.ID = chi.URLParam(r, "task_time_entry_id")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")

	err := handler.service.Update(r.Context(), projectID, taskID, taskTimeEntry)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update task time entry: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeEntry), nil)
}

func (handler *TaskTimerEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	taskTimeEntryID := chi.URLParam(r, "task_time_entry_id")
	err := handler.service.Delete(r.Context(), projectID, taskID, taskTimeEntryID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete task time entry: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *TaskTimerEntryHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	notes, err := handler.service.GetTaskTimerEntries(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}
