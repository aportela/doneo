package models

type TaskBase struct {
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type Task struct {
	ID          int          `json:"id"`
	Slug        string       `json:"slug"`
	Summary     string       `json:"summary"`
	Description string       `json:"description"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	CreatedBy   User         `json:"createdBy"`
	CreatedAt   int64        `json:"createdAt"`
	StartedAt   int64        `json:"startedAt"`
	FinishedAt  int64        `json:"finishedAt"`
	DueAt       int64        `json:"dueAt"`
	Tags        []string     `json:"tags"`
	Notes       []Note       `json:"notes"`
	Attachments []Attachment `json:"attachments"`
	LinkedTasks []TaskBase   `json:"linkedTasks"`
}
