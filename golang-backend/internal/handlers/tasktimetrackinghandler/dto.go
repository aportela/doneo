package tasktimetrackinghandler

import "github.com/aportela/doneo/internal/handlers/userhandler"

type addRequest struct {
	Summary   string `json:"summary"`
	SpentTime uint64 `json:"spentTime"`
}

type updateRequest struct {
	ID        string `json:"id"`
	Summary   string `json:"summary"`
	SpentTime uint64 `json:"spentTime"`
}

type taskTimeTrakingResponse struct {
	ID        string                       `json:"id"`
	CreatedBy userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt int64                        `json:"createdAt"`
	Summary   string                       `json:"summary"`
	SpentTime uint64                       `json:"spentTime"`
}

type searchResponse struct {
	TimeTrackings []taskTimeTrakingResponse `json:"timeTrackings"`
}
