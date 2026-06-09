package taskstatushandler

import "github.com/aportela/doneo/internal/handlers"

type statusFlags struct {
	DefaultStatusOnCreation bool `json:"defaultStatusOnCreation"`
	FillEmptyStartDate      bool `json:"fillEmptyStartDate"`
	SetStartDate            bool `json:"setStartDate"`
	FillEmptyFinishDate     bool `json:"fillEmptyFinishDate"`
	SetFinishDate           bool `json:"setFinishDate"`
}

type addRequest struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	HexColor string      `json:"hexColor"`
	Index    uint        `json:"index"`
	Flags    statusFlags `json:"flags"`
}

type updateRequest struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	HexColor string      `json:"hexColor"`
	Index    uint        `json:"index"`
	Flags    statusFlags `json:"flags"`
}

type filterRequest struct {
	Name *string `json:"name"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *filterRequest        `json:"filter"`
}

type TaskStatusResponse struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	HexColor string      `json:"hexColor"`
	Index    uint        `json:"index"`
	Flags    statusFlags `json:"flags"`
}

type searchResponse struct {
	TaskStatuses []TaskStatusResponse   `json:"taskStatuses"`
	Pager        handlers.PagerResponse `json:"pager"`
}
