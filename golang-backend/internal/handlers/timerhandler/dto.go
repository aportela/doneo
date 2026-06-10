package timerhandler

type TimerResponse struct {
	ID         string `json:"id"`
	StartedAt  int64  `json:"startedAt"`
	FinishedAt *int64 `json:"finishedAt"`
}

type searchResponse struct {
	Timers []TimerResponse `json:"timers"`
}
