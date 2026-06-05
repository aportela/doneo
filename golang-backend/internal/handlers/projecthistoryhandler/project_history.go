package projecthistoryhandler

import (
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projecthistoryservice"
	"github.com/go-chi/chi/v5"
)

type ProjectHistoryHandler struct {
	service projecthistoryservice.ProjectHistoryService
}

func NewHandler(service projecthistoryservice.ProjectHistoryService) *ProjectHistoryHandler {
	return &ProjectHistoryHandler{service: service}
}

func (handler *ProjectHistoryHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	projectAttachments, err := handler.service.Search(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
