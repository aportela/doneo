package timerhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/timerservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TimerHandler struct {
	service timerservice.TimerService
}

func NewHandler(service timerservice.TimerService) *TimerHandler {
	return &TimerHandler{service: service}
}

func (handler *TimerHandler) Start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request startRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserHandler] invalid request payload: %w", err))
		return
	}
	err := handler.service.Start(r.Context(), request.Summary)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TimerHandler] failed to start timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusCreated, handlers.ToEmptyResponse())
}

func (handler *TimerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timerId := chi.URLParam(r, "id")
	err := handler.service.Stop(r.Context(), timerId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TimerHandler] failed to stop timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *TimerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timerId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), timerId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TimerHandler] failed to delete timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *TimerHandler) Clear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := handler.service.Clear(r.Context())
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[TimerHandler] failed to clear timers: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *TimerHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timers, err := handler.service.Search(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(timers), err)
}
