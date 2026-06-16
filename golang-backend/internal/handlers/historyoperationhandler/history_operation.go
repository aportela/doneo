package historyoperationhandler

import (
	"net/http"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/go-chi/chi/v5"
)

type HistoryOperationHandler struct {
	db      database.Database
	service historyoperationservice.HistoryOperationService
}

func NewHandler(db database.Database, service historyoperationservice.HistoryOperationService) *HistoryOperationHandler {
	return &HistoryOperationHandler{db: db, service: service}
}

func (handler *HistoryOperationHandler) SearchProjectHistoryOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	historyOperations, err := handler.service.GetProjectHistoryOperations(r.Context(), handler.db, projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(historyOperations), err)
}

func (handler *HistoryOperationHandler) SearchTaskHistoryOperations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	historyOperations, err := handler.service.GetTaskHistoryOperations(r.Context(), handler.db, projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(historyOperations), err)
}
