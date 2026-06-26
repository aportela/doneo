package profilehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/profileservice"
)

type ProfileHandler interface {
	Update(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
}

type profileHandler struct {
	service profileservice.ProfileService
}

func NewHandler(service profileservice.ProfileService) ProfileHandler {
	return &profileHandler{service: service}
}

func (handler *profileHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProfileHandler] invalid request payload: %w", err))
		return
	}
	if user, err := handler.service.Update(r.Context(), updateRequestToDomain(request), request.Password); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProfileHandler] failed to update profile: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil)
	}
}

func (handler *profileHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if user, err := handler.service.Get(r.Context()); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProfileHandler] profile not found: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProfileHandler] failed to get profile: %w", err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, domainToResponse(user), nil)
	}
}
