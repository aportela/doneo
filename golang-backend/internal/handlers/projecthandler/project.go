package projecthandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	service projectservice.ProjectService
}

func NewHandler(service projectservice.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (handler *ProjectHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := addRequestToDomain(request)
	project, err := handler.service.Add(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to add project with ID %s: %w", request.ID, err))
		return
	}
	project, err = handler.service.Get(r.Context(), project.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get new project with ID %s: %w", project.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil, http.StatusCreated)
}

func (handler *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := updateRequestToDomain(request)
	project.ID = chi.URLParam(r, "id")
	project.UpdatedAt = utils.NowToTimePtr()
	project, err := handler.service.Update(r.Context(), project)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to update project with ID %s: %w", project.ID, err))
		return
	}
	project, err = handler.service.Get(r.Context(), project.ID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get updated project with ID %s: %w", request.ID, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (handler *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	err := handler.service.Delete(r.Context(), projectId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to delete project with ID %s: %w", projectId, err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	project, err := handler.service.Get(r.Context(), projectId)
	if err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] not found project with ID %s: %w", projectId, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get project with ID %s: %w", projectId, err))
			return
		}
	}
	handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
}

func (handler *ProjectHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectFilter{
		Key: nil,
	}
	if request.Filter != nil {
		if request.Filter.Key != nil {
			filter.Key = request.Filter.Key
		}
		if request.Filter.Summary != nil {
			filter.Summary = request.Filter.Summary
		}
		if request.Filter.TypeId != nil {
			filter.TypeId = request.Filter.TypeId
		}
		if request.Filter.PriorityId != nil {
			filter.PriorityId = request.Filter.PriorityId
		}
		if request.Filter.StatusId != nil {
			filter.StatusId = request.Filter.StatusId
		}
		if request.Filter.CreatedAt != nil {
			filter.CreatedAt = &domain.TimestampFilter{From: nil, To: nil}
			if request.Filter.CreatedAt.From != nil {
				filter.CreatedAt.From = request.Filter.CreatedAt.From
			}
			if request.Filter.CreatedAt.To != nil {
				filter.CreatedAt.To = request.Filter.CreatedAt.To
			}
		}
		if request.Filter.CreatedByUserId != nil {
			filter.CreatedByUserId = request.Filter.CreatedByUserId
		}
	}
	projects, pagerResult, err := handler.service.Search(r.Context(),
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
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projects, pagerResult), err)
}
