package notehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/go-chi/chi/v5"
)

type NoteHandler struct {
	service noteservice.NoteService
}

func NewHandler(service noteservice.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

func (handler *NoteHandler) AddProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")

	note, err := handler.service.AddProjectNote(r.Context(), projectID, note)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add note: %w", err))
		return
	}

	handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil, http.StatusCreated)
}

func (handler *NoteHandler) UpdateProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := updateRequestToDomain(request)
	note.ID = chi.URLParam(r, "note_id")
	projectID := chi.URLParam(r, "project_id")

	note, err := handler.service.UpdateProjectNote(r.Context(), projectID, note)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil)
}

func (handler *NoteHandler) DeleteProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	noteID := chi.URLParam(r, "note_id")
	err := handler.service.DeleteProjectNote(r.Context(), projectID, noteID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *NoteHandler) GetProjectNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "project_id")
	notes, err := handler.service.GetProjectNotes(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}

func (handler *NoteHandler) AddTaskNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")

	note, err := handler.service.AddTaskNote(r.Context(), projectID, taskID, note)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil, http.StatusCreated)
}

func (handler *NoteHandler) UpdateTaskNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := updateRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	note.ID = chi.URLParam(r, "note_id")

	note, err := handler.service.UpdateTaskNote(r.Context(), projectID, taskID, note)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil)
}

func (handler *NoteHandler) DeleteTaskNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	noteID := chi.URLParam(r, "note_id")
	err := handler.service.DeleteTaskNote(r.Context(), projectID, taskID, noteID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *NoteHandler) GetTaskNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	notes, err := handler.service.GetTaskNotes(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}
