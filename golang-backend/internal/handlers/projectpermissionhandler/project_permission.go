package projectpermissionhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectPermissionHandler struct {
	service projectpermissionservice.ProjectPermissionService
}

func NewHandler(db database.Database) *ProjectPermissionHandler {
	// TODO: rename to repository/service
	projectPermissionRepository := projectpermissionrepository.NewRepository(db)
	projectPermissionService := projectpermissionservice.NewService(projectPermissionRepository)
	return &ProjectPermissionHandler{service: projectPermissionService}
}

func (handler *ProjectPermissionHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] invalid request payload: %w", err))
		return
	}
	projectPermission := addRequestToDomain(request)
	projectPermission.ID = utils.UUID()
	projectId := chi.URLParam(r, "id")

	// TODO: fill User.Name && Role.Name for added notification on client
	err := handler.service.Add(r.Context(), projectPermission.ID, projectId, projectPermission.User.ID, projectPermission.Role.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] failed to add project permission: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, domainToResponse(projectPermission), nil, http.StatusCreated)
}

func (handler *ProjectPermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	permissionId := chi.URLParam(r, "permission_id")
	err := handler.service.Delete(r.Context(), projectId, permissionId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectPermissionHandler] failed to delete project permission: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *ProjectPermissionHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	projectPermissions, err := handler.service.Search(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectPermissions), err)
}
