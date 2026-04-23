package models

type ProjectType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Project struct {
	ID             int         `json:"id"`
	Key            string      `json:"key"`
	Summary        string      `json:"summary"`
	Description    string      `json:"description"`
	CreatedBy      User        `json:"createdBy"`
	CreatedAt      int64       `json:"createdAt"`
	LastModifiedAt int64       `json:"lastModifiedAt"`
	StartedAt      int64       `json:"startedAt"`
	FinishedAt     int64       `json:"finishedAt"`
	DueAt          int64       `json:"dueAt"`
	Type           ProjectType `json:"type"`
	//lead, asignee, priority, status
}
