package projectstatushandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/go-chi/chi/v5"
)

type ProjectStatusHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	SearchBase(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type projectStatusHandler struct {
	service projectstatusservice.ProjectStatusService
}

func NewHandler(service projectstatusservice.ProjectStatusService) ProjectStatusHandler {
	return &projectStatusHandler{service: service}
}

func (handler *projectStatusHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] invalid request payload: %w", err))
		return
	}
	projectStatus := addRequestToDomain(request)
	if projectStatus, err := handler.service.Add(r.Context(), projectStatus); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to add project status: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil, http.StatusCreated)
	}
}

func (handler *projectStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] invalid request payload: %w", err))
		return
	}
	projectStatus := updateRequestToDomain(request)
	projectStatus.ID = chi.URLParam(r, "id")
	if projectStatus, err := handler.service.Update(r.Context(), projectStatus); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to update project status: %w", err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil)
	}
}

func (handler *projectStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectStatusID := chi.URLParam(r, "id")
	if err := handler.service.Delete(r.Context(), projectStatusID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to delete project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *projectStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectStatusID := chi.URLParam(r, "id")
	if projectStatus, err := handler.service.Get(r.Context(), projectStatusID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to get non existent project status: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to get projectStatus: %w", err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil)
	}
}

func (handler *projectStatusHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectStatuses, pagerResult, err := handler.service.SearchBase(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectStatuses, pagerResult), err)
}

func (handler *projectStatusHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectStatusesFilter{
		Name: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
	}
	projectStatuses, pagerResult, err := handler.service.Search(r.Context(),
		browser.Params{
			CurrentPage: request.Pager.CurrentPage,
			ResultsPage: request.Pager.ResultsPage,
		},
		browser.Order{
			Field: request.Order.Field,
			Sort:  string(request.Order.Sort),
		},
		filter,
	)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectStatuses, pagerResult), err)
}
