package historyoperationhandler

import (
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/go-chi/chi/v5"
)

type HistoryOperationHandler struct {
	service historyoperationservice.HistoryOperationService
}

func NewHandler(service historyoperationservice.HistoryOperationService) *HistoryOperationHandler {
	return &HistoryOperationHandler{service: service}
}

func (handler *HistoryOperationHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	historyOperations, err := handler.service.Search(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(historyOperations), err)
}
