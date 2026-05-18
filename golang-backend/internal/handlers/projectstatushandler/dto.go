package projectstatushandler

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

type projectStatusResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	HexColor string `json:"hexColor"`
	Index    int    `json:"index"`
}

type searchResponse struct {
	ProjectStatuses []projectStatusResponse `json:"projectStatuses"`
	Pager           handlers.PagerResponse  `json:"pager"`
}
