package taskhandler

import (
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type taskPriority struct {
	ID string `json:"id"`
}

type taskStatus struct {
	ID string `json:"id"`
}

type addRequest struct {
	ID          string       `json:"id"`
	Summary     string       `json:"summary"`
	Description *string      `json:"description"`
	Priority    taskPriority `json:"priority"`
	Status      taskStatus   `json:"status"`
	Tags        []string     `json:"tags"`
}

type updateRequest struct {
	ID          string       `json:"id"`
	Summary     string       `json:"summary"`
	Description *string      `json:"description"`
	Priority    taskPriority `json:"priority"`
	Status      taskStatus   `json:"status"`
	StartedAt   *int64       `json:"startedAt"`
	FinishedAt  *int64       `json:"finishedAt"`
	DueAt       *int64       `json:"dueAt"`
	Tags        []string     `json:"tags"`
}

type filterRequest struct {
	Summary         *string                   `json:"summary"`
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

type taskResponse struct {
	ID                     string                                   `json:"id"`
	ProjectID              string                                   `json:"projectId"`
	Slug                   string                                   `json:"slug"`
	Summary                string                                   `json:"summary"`
	Description            *string                                  `json:"description"`
	CreatedBy              userhandler.UserBaseResponse             `json:"createdBy"`
	CreatedAt              int64                                    `json:"createdAt"`
	UpdatedAt              *int64                                   `json:"updatedAt"`
	DeletedAt              *int64                                   `json:"deletedAt"`
	StartedAt              *int64                                   `json:"startedAt"`
	FinishedAt             *int64                                   `json:"finishedAt"`
	DueAt                  *int64                                   `json:"dueAt"`
	Priority               taskpriorityhandler.TaskPriorityResponse `json:"priority"`
	Status                 taskstatushandler.TaskStatusResponse     `json:"status"`
	Tags                   []string                                 `json:"tags"`
	TasksCount             uint                                     `json:"tasksCount"`
	PermissionsCount       uint                                     `json:"permissionsCount"`
	AttachmentsCount       uint                                     `json:"attachmentsCount"`
	NotesCount             uint                                     `json:"notesCount"`
	HistoryOperationsCount uint                                     `json:"historyOperationsCount"`
	TimeTrackingsCount     uint                                     `json:"timeTrackingsCount"`
}

type searchTasksResponse struct {
	Tasks []taskResponse         `json:"tasks"`
	Pager handlers.PagerResponse `json:"pager"`
}
