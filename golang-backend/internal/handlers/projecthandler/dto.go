package projecthandler

import (
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type projectType struct {
	ID string `json:"id"`
}

type projectPriority struct {
	ID string `json:"id"`
}

type projectStatus struct {
	ID string `json:"id"`
}

type addRequest struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Summary     string          `json:"summary"`
	Description *string         `json:"description"`
	Type        projectType     `json:"type"`
	Priority    projectPriority `json:"priority"`
	Status      projectStatus   `json:"status"`
}

type updateRequest struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Summary     string          `json:"summary"`
	Description *string         `json:"description"`
	Type        projectType     `json:"type"`
	Priority    projectPriority `json:"priority"`
	Status      projectStatus   `json:"status"`
	StartedAt   *int64          `json:"startedAt"`
	FinishedAt  *int64          `json:"finishedAt"`
	DueAt       *int64          `json:"dueAt"`
}

type filterRequest struct {
	Slug            *string                   `json:"slug"`
	Summary         *string                   `json:"summary"`
	TypeID          *string                   `json:"typeId"`
	PriorityID      *string                   `json:"priorityId"`
	StatusID        *string                   `json:"statusId"`
	CreatedAt       *handlers.TimestampFilter `json:"createdAt"`
	CreatedByUserID *string                   `json:"createdByUserId"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *filterRequest        `json:"filter"`
}

type projectResponse struct {
	ID                     string                                         `json:"id"`
	SLug                   string                                         `json:"slug"`
	Summary                string                                         `json:"summary"`
	Description            *string                                        `json:"description"`
	CreatedBy              userhandler.UserBaseResponse                   `json:"createdBy"`
	CreatedAt              int64                                          `json:"createdAt"`
	UpdatedAt              *int64                                         `json:"updatedAt"`
	DeletedAt              *int64                                         `json:"deletedAt"`
	StartedAt              *int64                                         `json:"startedAt"`
	FinishedAt             *int64                                         `json:"finishedAt"`
	DueAt                  *int64                                         `json:"dueAt"`
	Type                   projecttypehandler.ProjectTypeResponse         `json:"type"`
	Priority               projectpriorityhandler.ProjectPriorityResponse `json:"priority"`
	Status                 projectstatushandler.ProjectStatusResponse     `json:"status"`
	TasksCount             uint16                                         `json:"tasksCount"`
	PermissionsCount       uint16                                         `json:"permissionsCount"`
	AttachmentsCount       uint16                                         `json:"attachmentsCount"`
	NotesCount             uint16                                         `json:"notesCount"`
	HistoryOperationsCount uint16                                         `json:"historyOperationsCount"`
}

type searchProjectsResponse struct {
	Projects []projectResponse      `json:"projects"`
	Pager    handlers.PagerResponse `json:"pager"`
}
