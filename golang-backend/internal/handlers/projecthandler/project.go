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

type ProjectHandler interface {
	Add(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
}

type projectHandler struct {
	service projectservice.ProjectService
}

func NewHandler(service projectservice.ProjectService) ProjectHandler {
	return &projectHandler{service: service}
}

func (handler *projectHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := addRequestToDomain(request)
	if project, err := handler.service.Add(r.Context(), project); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to add project with ID %s: %w", request.ID, err))
		return
	} else {
		if project, err := handler.service.Get(r.Context(), project.ID); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get new project with ID %s: %w", project.ID, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil, http.StatusCreated)
		}
	}
}

func (handler *projectHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request updateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := updateRequestToDomain(request)
	project.ID = chi.URLParam(r, "project_id")
	project.UpdatedAt = utils.NowToTimePtr()
	if project, err := handler.service.Update(r.Context(), project); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to update project with ID %s: %w", project.ID, err))
		return
	} else {
		if project, err := handler.service.Get(r.Context(), project.ID); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get updated project with ID %s: %w", request.ID, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
		}
	}
}

func (handler *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	if err := handler.service.Delete(r.Context(), projectID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to delete project with ID %s: %w", projectID, err))
		return
	} else {
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
}

func (handler *projectHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	if project, err := handler.service.Get(r.Context(), projectID); err != nil {
		if err == domain.NotFoundError {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] not found project with ID %s: %w", projectID, err))
			return
		} else {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get project with ID %s: %w", projectID, err))
			return
		}
	} else {
		handlers.ToHandlerJSONResponse(w, DomainToResponse(project), nil)
	}
}

func (handler *projectHandler) Search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request searchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	filter := domain.SearchProjectFilter{
		Slug: nil,
	}
	if request.Filter != nil {
		if request.Filter.Slug != nil {
			filter.Slug = request.Filter.Slug
		}
		if request.Filter.Summary != nil {
			filter.Summary = request.Filter.Summary
		}
		if request.Filter.TypeID != nil {
			filter.TypeID = request.Filter.TypeID
		}
		if request.Filter.PriorityID != nil {
			filter.PriorityID = request.Filter.PriorityID
		}
		if request.Filter.StatusID != nil {
			filter.StatusID = request.Filter.StatusID
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
		if request.Filter.CreatedByUserID != nil {
			filter.CreatedByUserID = request.Filter.CreatedByUserID
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
