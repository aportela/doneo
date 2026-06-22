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

type UserTimerHandler interface {
	StartUserTimer(w http.ResponseWriter, r *http.Request)
	StopUserTimer(w http.ResponseWriter, r *http.Request)
	DeleteUserTimer(w http.ResponseWriter, r *http.Request)
	ClearUserTimers(w http.ResponseWriter, r *http.Request)
	GetUserTimers(w http.ResponseWriter, r *http.Request)
}

type userTimerHandler struct {
	userTimerService usertimerservice.UserTimerService
}

func NewHandler(userTimerService usertimerservice.UserTimerService) UserTimerHandler {
	return &userTimerHandler{userTimerService: userTimerService}
}

func (handler *userTimerHandler) StartUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request startRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] invalid request payload: %w", err))
		return
	}
	if err := handler.userTimerService.StartUserTimer(r.Context(), request.Summary); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to start user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusCreated, handlers.ToEmptyResponse())
}

func (handler *userTimerHandler) StopUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userTimerID := chi.URLParam(r, "id")
	if err := handler.userTimerService.StopUserTimer(r.Context(), userTimerID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to stop user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *userTimerHandler) DeleteUserTimer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userTimerID := chi.URLParam(r, "id")
	if err := handler.userTimerService.DeleteUserTimer(r.Context(), userTimerID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to delete user timer: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *userTimerHandler) ClearUserTimers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := handler.userTimerService.ClearUserTimers(r.Context()); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[UserTimerHandler] failed to clear user timers: %w", err))
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, handlers.ToEmptyResponse())
}

func (handler *userTimerHandler) GetUserTimers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	timers, err := handler.userTimerService.GetUserTimers(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(timers), err)
}
