package timerhandler

type startRequest struct {
	Summary string `json:"summary"`
}

type timerResponse struct {
	ID         string `json:"id"`
	Summary    string `json:"summary"`
	StartedAt  int64  `json:"startedAt"`
	FinishedAt *int64 `json:"finishedAt"`
}

type searchResponse struct {
	Timers []timerResponse `json:"timers"`
}
