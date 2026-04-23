package models

type TaskBase struct {
	ID      int    `json:"id"`
	Summary string `json:"summary"`
}
type Task struct {
	ID          int          `json:"id"`
	Summary     string       `json:"summary"`
	Status      Status       `json:"status"`
	Priority    Priority     `json:"priority"`
	CreatedBy   User         `json:"createdBy"`
	CreatedAt   int64        `json:"createdAt"`
	FinishedAt  int64        `json:"finishedAt"`
	StartDate   int64        `json:"startDate"`
	DueDate     int64        `json:"dueDate"`
	Tags        []string     `json:"tags"`
	Notes       []Note       `json:"notes"`
	Attachments []Attachment `json:"attachments"`
}
