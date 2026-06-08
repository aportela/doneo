package projectstatushandler

import "github.com/aportela/doneo/internal/handlers"

type addRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    uint   `json:"index"`
}

type updateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    uint   `json:"index"`
}

type filterRequest struct {
	Name *string `json:"name"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *filterRequest        `json:"filter"`
}

type ProjectStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    uint   `json:"index"`
}

type searchResponse struct {
	ProjectStatuses []ProjectStatusResponse `json:"projectStatuses"`
	Pager           handlers.PagerResponse  `json:"pager"`
}
