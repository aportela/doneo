package tasktimerentryhandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type addRequest struct {
	Summary      string `json:"summary"`
	TotalSeconds uint   `json:"totalSeconds"`
}

type updateRequest struct {
	ID           string `json:"id"`
	Summary      string `json:"summary"`
	TotalSeconds uint   `json:"totalSeconds"`
}

type TaskTimeEntryResponse struct {
	ID           string                       `json:"id"`
	User         userhandler.UserBaseResponse `json:"user"`
	CreatedAt    int64                        `json:"createdAt"`
	Summary      string                       `json:"summary"`
	TotalSeconds uint                         `json:"totalSeconds"`
}

type searchResponse struct {
	ProjectPermissions []TaskTimeEntryResponse `json:"timeEntries"`
}
