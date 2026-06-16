package usertimerhandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/usertimerservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserTimerHandler struct {
	userTimerService usertimerservice.UserTimerService
}

func NewHandler(userTimerService usertimerservice.UserTimerService) *UserTimerHandler {
	return &UserTimerHandler{userTimerService: userTimerService}
}

func (handler *UserTimerHandler) StartUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request startRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] invalid request payload: %w", err))
		return
	}
	err := handler.userTimerService.StartUserTimer(r.Context(), request.Summary)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to start user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusCreated, handlers.ToEmptyResponse())
}

func (handler *UserTimerHandler) StopUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userTimerID := chi.URLParam(r, "id")
	err := handler.userTimerService.StopUserTimer(r.Context(), userTimerID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to stop user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *UserTimerHandler) DeleteUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userTimerID := chi.URLParam(r, "id")
	err := handler.userTimerService.DeleteUserTimer(r.Context(), userTimerID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to delete user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *UserTimerHandler) ClearUserTimers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := handler.userTimerService.ClearUserTimers(r.Context())
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to clear user timers: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *UserTimerHandler) GetUserTimers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timers, err := handler.userTimerService.GetUserTimers(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(timers), err)
}
