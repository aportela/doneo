package tasktimetrackinghandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type addRequest struct {
	Summary      string `json:"summary"`
	TotalSeconds uint64 `json:"totalSeconds"`
}

type updateRequest struct {
	ID           string `json:"id"`
	Summary      string `json:"summary"`
	TotalSeconds uint64 `json:"totalSeconds"`
}

type taskTimeTrakingResponse struct {
	ID           string                       `json:"id"`
	CreatedBy    userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt    int64                        `json:"createdAt"`
	Summary      string                       `json:"summary"`
	TotalSeconds uint64                       `json:"totalSeconds"`
}

type searchResponse struct {
	TimeTrackings []taskTimeTrakingResponse `json:"timeTrackings"`
}
