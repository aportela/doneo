package projectstatushandler

import "github.com/aportela/doneo/internal/handlers"

type addRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type updateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type filterRequest struct {
	Name *string `json:"name"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *filterRequest        `json:"filter"`
}

type projectStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type searchResponse struct {
	ProjectStatuses []projectStatusResponse `json:"projectStatuses"`
	Pager           handlers.PagerResponse  `json:"pager"`
}
