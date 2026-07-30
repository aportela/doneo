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
	Patch(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Search(w http.ResponseWriter, r *http.Request)
	GetCurrentProjects(w http.ResponseWriter, r *http.Request)
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

func (handler *projectHandler) Patch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var request patchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] invalid request payload: %w", err))
		return
	}
	project := patchRequestToDomain(request)
	project.ID = chi.URLParam(r, "project_id")
	if project, err := handler.service.Patch(r.Context(), project); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to patch project with ID %s: %w", project.ID, err))
		return
	} else {
		if project, err := handler.service.Get(r.Context(), project.ID); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[ProjectHandler] failed to get patched project with ID %s: %w", request.ID, err))
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
		if request.Filter.CreatedByUserID != nil {
			filter.CreatedByUserID = request.Filter.CreatedByUserID
		}
		if request.Filter.CreatedAt != nil {
			filter.CreatedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.CreatedAt.Filled != nil {
				filter.CreatedAt.Filled = request.Filter.CreatedAt.Filled
			} else if request.Filter.CreatedAt.Empty != nil {
				filter.CreatedAt.Empty = request.Filter.CreatedAt.Empty
			} else {
				if request.Filter.CreatedAt.From != nil {
					filter.CreatedAt.From = request.Filter.CreatedAt.From
				}
				if request.Filter.CreatedAt.To != nil {
					filter.CreatedAt.To = request.Filter.CreatedAt.To
				}
			}
		}
		if request.Filter.UpdatedAt != nil {
			filter.UpdatedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.UpdatedAt.Filled != nil {
				filter.UpdatedAt.Filled = request.Filter.UpdatedAt.Filled
			} else if request.Filter.UpdatedAt.Empty != nil {
				filter.UpdatedAt.Empty = request.Filter.UpdatedAt.Empty
			} else {
				if request.Filter.UpdatedAt.From != nil {
					filter.UpdatedAt.From = request.Filter.UpdatedAt.From
				}
				if request.Filter.UpdatedAt.To != nil {
					filter.UpdatedAt.To = request.Filter.UpdatedAt.To
				}
			}
		}
		if request.Filter.StartedAt != nil {
			filter.StartedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.StartedAt.Filled != nil {
				filter.StartedAt.Filled = request.Filter.StartedAt.Filled
			} else if request.Filter.StartedAt.Empty != nil {
				filter.StartedAt.Empty = request.Filter.StartedAt.Empty
			} else {
				if request.Filter.StartedAt.From != nil {
					filter.StartedAt.From = request.Filter.StartedAt.From
				}
				if request.Filter.StartedAt.To != nil {
					filter.StartedAt.To = request.Filter.StartedAt.To
				}
			}
		}
		if request.Filter.FinishedAt != nil {
			filter.FinishedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.FinishedAt.Filled != nil {
				filter.FinishedAt.Filled = request.Filter.FinishedAt.Filled
			} else if request.Filter.FinishedAt.Empty != nil {
				filter.FinishedAt.Empty = request.Filter.FinishedAt.Empty
			} else {
				if request.Filter.FinishedAt.From != nil {
					filter.FinishedAt.From = request.Filter.FinishedAt.From
				}
				if request.Filter.FinishedAt.To != nil {
					filter.FinishedAt.To = request.Filter.FinishedAt.To
				}
			}
		}
		if request.Filter.DueAt != nil {
			filter.DueAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.DueAt.Filled != nil {
				filter.DueAt.Filled = request.Filter.DueAt.Filled
			} else if request.Filter.DueAt.Empty != nil {
				filter.DueAt.Empty = request.Filter.DueAt.Empty
			} else {
				if request.Filter.DueAt.From != nil {
					filter.DueAt.From = request.Filter.DueAt.From
				}
				if request.Filter.DueAt.To != nil {
					filter.DueAt.To = request.Filter.DueAt.To
				}
			}
		}
		if request.Filter.ArchivedAt != nil {
			filter.ArchivedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.ArchivedAt.Filled != nil {
				filter.ArchivedAt.Filled = request.Filter.ArchivedAt.Filled
			} else if request.Filter.ArchivedAt.Empty != nil {
				filter.ArchivedAt.Empty = request.Filter.ArchivedAt.Empty
			} else {
				if request.Filter.ArchivedAt.From != nil {
					filter.ArchivedAt.From = request.Filter.ArchivedAt.From
				}
				if request.Filter.ArchivedAt.To != nil {
					filter.ArchivedAt.To = request.Filter.ArchivedAt.To
				}
			}
		}
		if request.Filter.DeletedAt != nil {
			filter.DeletedAt = &domain.TimestampFilter{From: nil, To: nil, Filled: nil, Empty: nil}
			if request.Filter.DeletedAt.Filled != nil {
				filter.DeletedAt.Filled = request.Filter.DeletedAt.Filled
			} else if request.Filter.DeletedAt.Empty != nil {
				filter.DeletedAt.Empty = request.Filter.DeletedAt.Empty
			} else {
				if request.Filter.DeletedAt.From != nil {
					filter.DeletedAt.From = request.Filter.DeletedAt.From
				}
				if request.Filter.DeletedAt.To != nil {
					filter.DeletedAt.To = request.Filter.DeletedAt.To
				}
			}
		}
	}
	projects, pagerResult, err := handler.service.Search(r.Context(),
		browser.PagerQuery{
			Enabled:     request.Pager.Enabled,
			CurrentPage: request.Pager.CurrentPage,
			ResultsPage: request.Pager.ResultsPage,
		},
		browser.Order{
			Field:     request.Order.Field,
			Direction: string(request.Order.Direction),
		},
		filter,
	)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projects, pagerResult), err)
}

func (handler *projectHandler) GetCurrentProjects(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projects, err := handler.service.GetCurrentProjects(r.Context())
	handlers.ToHandlerJSONResponse(w, map[string]any{"projects": domainArrayToResponseArray(projects)}, err)
}
