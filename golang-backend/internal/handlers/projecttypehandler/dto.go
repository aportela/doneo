package projecttypehandler

import "github.com/aportela/doneo/internal/handlers"

type addRequest struct {
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type updateRequest struct {
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

type ProjectTypeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type searchResponse struct {
	ProjectTypes []ProjectTypeResponse  `json:"projectTypes"`
	Pager        handlers.PagerResponse `json:"pager"`
}
