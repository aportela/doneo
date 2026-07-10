package pagehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/pageservice"
	"github.com/go-chi/chi/v5"
)

type PageHandler interface {
	AddProjectPage(w http.ResponseWriter, r *http.Request)
	UpdateProjectPage(w http.ResponseWriter, r *http.Request)
	DeleteProjectPage(w http.ResponseWriter, r *http.Request)
	GetProjectPage(w http.ResponseWriter, r *http.Request)
	GetProjectPages(w http.ResponseWriter, r *http.Request)
}

type pageHandler struct {
	service pageservice.PageService
}

func NewHandler(service pageservice.PageService) PageHandler {
	return &pageHandler{service: service}
}

func (handler *pageHandler) AddProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[PageHandler] invalid request payload: %w", err))
		return
	}
	page := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")
	if page, err := handler.service.AddProjectPage(r.Context(), projectID, page); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[PageHandler] failed to add page: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(page), nil, http.StatusCreated)
	}
}

func (handler *pageHandler) UpdateProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[PageHandler] invalid request payload: %w", err))
		return
	}
	page := updateRequestToDomain(request)
	page.ID = chi.URLParam(r, "page_id")
	projectID := chi.URLParam(r, "project_id")
	if page, err := handler.service.UpdateProjectPage(r.Context(), projectID, page); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[PageHandler] failed to update page: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(page), nil)
	}
}

func (handler *pageHandler) DeleteProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	pageID := chi.URLParam(r, "page_id")
	if err := handler.service.DeleteProjectPage(r.Context(), projectID, pageID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[PageHandler] failed to delete page: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
}

func (handler *pageHandler) GetProjectPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	pageID := chi.URLParam(r, "page_id")
	page, err := handler.service.GetProjectPage(r.Context(), projectID, pageID)
	handlers.ToHandlerJSONResponse(w, domainToResponse(page), err)
}

func (handler *pageHandler) GetProjectPages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	pages, err := handler.service.GetProjectPages(r.Context(), projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(pages), err)
}
