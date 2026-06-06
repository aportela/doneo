package notehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/noteservice"
	"github.com/aportela/doneo/internal/utils"
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
	note.ID = utils.UUID()
	projectId := chi.URLParam(r, "id")

	err := handler.service.AddProjectNote(r.Context(), projectId, note)
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

	projectId := chi.URLParam(r, "id")
	// TODO: move to service ????
	note.UpdatedAt = utils.NowToTimePtr()

	err := handler.service.UpdateProjectNote(r.Context(), projectId, note)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to add note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(note), nil)
}

func (handler *NoteHandler) DeleteProjectNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	noteId := chi.URLParam(r, "note_id")
	err := handler.service.DeleteProjectNote(r.Context(), projectId, noteId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] failed to delete note: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *NoteHandler) GetProjectNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	projectPermissions, err := handler.service.GetProjectNotes(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectPermissions), err)
}
