package projecthistoryhandler

import (
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrespository"
	"github.com/aportela/doneo/internal/services/projecthistoryservice"
	"github.com/go-chi/chi/v5"
)

type ProjectHistoryHandler struct {
	service  projecthistoryservice.ProjectHistoryService
	basePath string
}

func NewHandler(db database.Database, basePath string) *ProjectHistoryHandler {
	repository := projecthistoryrespository.NewRepository(db)
	service := projecthistoryservice.NewService(repository)
	return &ProjectHistoryHandler{service: service, basePath: basePath}
}

func (handler *ProjectHistoryHandler) GetProjectHistoryOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	projectAttachments, err := handler.service.GetProjectHistoryOperations(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
