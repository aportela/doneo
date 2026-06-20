package projecttypehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/go-chi/chi/v5"
)

type ProjectTypeHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	SearchBase(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type projectTypeHandler struct {
	service projecttypeservice.ProjectTypeService
}

func NewHandler(service projecttypeservice.ProjectTypeService) ProjectTypeHandler {
	return &projectTypeHandler{service: service}
}

func (handler *projectTypeHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	projectType := addRequestToDomain(request)
	projectType, err := handler.service.Add(r.Context(), projectType)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to add project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectType), nil, http.StatusCreated)
}

func (handler *projectTypeHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	projectType := updateRequestToDomain(request)
	projectType.ID = chi.URLParam(r, "id")
	projectType, err := handler.service.Update(r.Context(), projectType)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to update project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectType), nil)
}

func (handler *projectTypeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypeID := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), projectTypeID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to delete project type: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *projectTypeHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypeID := chi.URLParam(r, "id")
	projectType, err := handler.service.Get(r.Context(), projectTypeID)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to get non existent project type: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] failed to get project type: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectType), nil)
}

func (handler *projectTypeHandler) SearchBase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectTypes, pagerResult, err := handler.service.SearchBase(r.Context())
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectTypes, pagerResult), err)
}

func (handler *projectTypeHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectTypeHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectTypesFilter{
		Name: nil,
	}
	if request.Filter != nil {
		if request.Filter.Name != nil {
			filter.Name = request.Filter.Name
		}
	}
	projectTypes, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectTypes, pagerResult), err)
}
