package notehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/go-chi/chi/v5"
)

type NoteHandler interface {
	AddProjectNote(w http.ResponseWriter, r *http.Request)
	UpdateProjectNote(w http.ResponseWriter, r *http.Request)
	DeleteProjectNote(w http.ResponseWriter, r *http.Request)
	GetProjectNotes(w http.ResponseWriter, r *http.Request)
	AddTaskNote(w http.ResponseWriter, r *http.Request)
	UpdateTaskNote(w http.ResponseWriter, r *http.Request)
	DeleteTaskNote(w http.ResponseWriter, r *http.Request)
	GetTaskNotes(w http.ResponseWriter, r *http.Request)
}

type noteHandler struct {
	service noteservice.NoteService
}

func NewHandler(service noteservice.NoteService) NoteHandler {
	return &noteHandler{service: service}
}

func (handler *noteHandler) AddProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	if note, err := handler.service.AddProjectNote(r.Context(), projectID, note); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add note: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil, http.StatusCreated)
	}
}

func (handler *noteHandler) UpdateProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := updateRequestToDomain(request)
	note.ID = chi.URLParam(r, "note_id")
	projectID := chi.URLParam(r, "project_id")
	if note, err := handler.service.UpdateProjectNote(r.Context(), projectID, note); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update note: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil)
	}
}

func (handler *noteHandler) DeleteProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	noteID := chi.URLParam(r, "note_id")
	if err := handler.service.DeleteProjectNote(r.Context(), projectID, noteID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete note: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
}

func (handler *noteHandler) GetProjectNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	notes, err := handler.service.GetProjectNotes(r.Context(), projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}

func (handler *noteHandler) AddTaskNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	note := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	if note, err := handler.service.AddTaskNote(r.Context(), projectID, taskID, note); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add note: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil, http.StatusCreated)
	}
}

func (handler *noteHandler) UpdateTaskNote(w http.ResponseWriter, r *http.Request) {
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
	if note, err := handler.service.UpdateTaskNote(r.Context(), projectID, taskID, note); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to update note: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil)
	}
}

func (handler *noteHandler) DeleteTaskNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	noteID := chi.URLParam(r, "note_id")
	if err := handler.service.DeleteTaskNote(r.Context(), projectID, taskID, noteID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *noteHandler) GetTaskNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	notes, err := handler.service.GetTaskNotes(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(notes), err)
}
