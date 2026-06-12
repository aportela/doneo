package taskhistoryhandler

import (
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/taskhistoryservice"
	"github.com/go-chi/chi/v5"
)

type TaskHistoryHandler struct {
	service taskhistoryservice.TaskHistoryService
}

func NewHandler(service taskhistoryservice.TaskHistoryService) *TaskHistoryHandler {
	return &TaskHistoryHandler{service: service}
}

func (handler *TaskHistoryHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	taskId := chi.URLParam(r, "task_id")
	historyOperations, err := handler.service.Search(r.Context(), taskId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(historyOperations), err)
}
