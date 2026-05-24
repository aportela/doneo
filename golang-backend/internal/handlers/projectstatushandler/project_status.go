package projectstatushandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectStatusHandler struct {
	service projectstatusservice.ProjectStatusService
}

func NewProjectStatusHandler(db database.Database) *ProjectStatusHandler {
	projectStatusRepository := projectstatusrepository.NewProjectStatusRepository(db)
	projectStatusService := projectstatusservice.NewProjectStatusService(projectStatusRepository)
	return &ProjectStatusHandler{service: projectStatusService}
}

func (h *ProjectStatusHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] invalid request payload: %w", err))
		return
	}
	projectStatus := addRequestToDomain(request)
	projectStatus.ID = utils.UUID()
	err := h.service.Add(r.Context(), projectStatus)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to add project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil, http.StatusCreated)
}

func (h *ProjectStatusHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] invalid request payload: %w", err))
		return
	}
	projectStatus := updateRequestToDomain(request)
	projectStatus.ID = chi.URLParam(r, "id")
	err := h.service.Update(r.Context(), projectStatus)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to update project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil)
}

func (h *ProjectStatusHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectStatusId := chi.URLParam(r, "id")
	err := h.service.Delete(r.Context(), projectStatusId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to delete project status: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (h *ProjectStatusHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectStatusId := chi.URLParam(r, "id")
	projectStatus, err := h.service.Get(r.Context(), projectStatusId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to get non existent project status: %w", err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectStatusHandler] failed to get projectStatus: %w", err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(projectStatus), nil)
}

func (h *ProjectStatusHandler) Search(w http.ResponseWriter, r *http.Request) {
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
	projectStatuses, pagerResult, err := h.service.Search(r.Context(),
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
