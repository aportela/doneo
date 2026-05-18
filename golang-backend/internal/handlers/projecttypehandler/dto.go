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

type projectTypeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
}

type searchResponse struct {
	ProjectTypes []projectTypeResponse  `json:"projectTypes"`
	Pager        handlers.PagerResponse `json:"pager"`
}
