package models

type ProjectType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID          int         `json:"id"`
	Key         string      `json:"key"`
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	CreatedBy   User        `json:"createdBy"`
	CreatedAt   int64       `json:"createdAt"`
	FinishedAt  int64       `json:"finishedAt"`
	StartDate   int64       `json:"startDate"`
	DueDate     int64       `json:"dueDate"`
	Type        ProjectType `json:"type"`
	Tags        []string    `json:"tags"`
	//lead, asignee, priority, status
}
