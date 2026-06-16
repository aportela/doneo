package tasktimeentryhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/tasktimeentryservice"
	"github.com/go-chi/chi/v5"
)

type TaskTimeEntryHandler struct {
	service tasktimeentryservice.TaskTimerEntryService
}

func NewHandler(service tasktimeentryservice.TaskTimerEntryService) *TaskTimeEntryHandler {
	return &TaskTimeEntryHandler{service: service}
}

func (handler *TaskTimeEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	taskTimeEntry := addRequestToDomain(request)
	projectId := chi.URLParam(r, "project_id")
	taskId := chi.URLParam(r, "task_id")

	// TODO: return taskTimeEntry with new id & createdAt
	err := handler.service.Add(r.Context(), projectId, taskId, taskTimeEntry)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add task time entry: %w", err))
		return
	}

	handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeEntry), nil, http.StatusCreated)
}

func (handler *TaskTimeEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	taskTimeEntry := updateRequestToDomain(request)
	taskTimeEntry.ID = chi.URLParam(r, "task_time_entry_id")
	projectId := chi.URLParam(r, "project_id")
	taskId := chi.URLParam(r, "task_id")

	err := handler.service.Update(r.Context(), projectId, taskId, taskTimeEntry)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update task time entry: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(taskTimeEntry), nil)
}

func (handler *TaskTimeEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "project_id")
	taskId := chi.URLParam(r, "task_id")
	taskTimeEntryId := chi.URLParam(r, "task_time_entry_id")
	err := handler.service.Delete(r.Context(), projectId, taskId, taskTimeEntryId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete task time entry: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *TaskTimeEntryHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskId := chi.URLParam(r, "task_id")
	notes, err := handler.service.GetTaskTimeEntries(r.Context(), taskId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}
