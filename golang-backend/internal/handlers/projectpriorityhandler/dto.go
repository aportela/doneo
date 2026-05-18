package projectpriorityhandler

import "github.com/aportela/doneo/internal/handlers"

type addRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    int    `json:"index"`
}

type updateRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    int    `json:"index"`
}

type projectPriorityResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    int    `json:"index"`
}

type searchResponse struct {
	ProjectPriorities []projectPriorityResponse `json:"projectPriorities"`
	Pager             handlers.PagerResponse    `json:"pager"`
}
