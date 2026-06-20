package projectpermissionhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/go-chi/chi/v5"
)

type ProjectPermissionHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetProjectPermissions(w http.ResponseWriter, r *http.Request)
}

type projectPermissionHandler struct {
	service projectpermissionservice.ProjectPermissionService
}

func NewHandler(service projectpermissionservice.ProjectPermissionService) ProjectPermissionHandler {
	return &projectPermissionHandler{service: service}
}

func (handler *projectPermissionHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] invalid request payload: %w", err))
		return
	}
	projectPermission := addRequestToDomain(request)
	projectID := chi.URLParam(r, "project_id")

	projectPermission, err := handler.service.Add(r.Context(), projectID, projectPermission)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] failed to add project permission: %w", err))
		return
	}

	handlers.ToHandlerJSONResponse(w, domainToResponse(projectPermission), nil, http.StatusCreated)
}

func (handler *projectPermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	projectPermissionID := chi.URLParam(r, "permission_id")
	err := handler.service.Delete(r.Context(), projectID, projectPermissionID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] failed to delete project permission: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *projectPermissionHandler) GetProjectPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	projectPermissions, err := handler.service.GetProjectPermissions(r.Context(), projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectPermissions), err)
}
