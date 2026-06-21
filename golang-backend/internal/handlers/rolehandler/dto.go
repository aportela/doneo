package rolehandler

import "github.com/aportela/doneo/internal/handlers"

type permissionsFlags struct {
	AllowUpdateProject bool `json:"allowUpdateProject"`
	AllowDeleteProject bool `json:"allowDeleteProject"`
	AllowViewProject   bool `json:"allowViewProject"`
	AllowAddTask       bool `json:"allowAddTask"`
	AllowUpdateTask    bool `json:"allowUpdateTask"`
	AllowDeleteTask    bool `json:"allowDeleteTask"`
	AllowViewTask      bool `json:"allowViewTask"`
}

type addRequest struct {
	Name        string           `json:"name"`
	Permissions permissionsFlags `json:"permissions"`
}

type updateRequest struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions permissionsFlags `json:"permissions"`
}

type filterRequest struct {
	Name *string `json:"name"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *filterRequest        `json:"filter"`
}

type RoleBaseResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoleResponse struct {
	RoleBaseResponse
	Permissions permissionsFlags `json:"permissions"`
}

type searchBaseResponse struct {
	Roles []RoleBaseResponse `json:"roles"`
}

type searchResponse struct {
	Roles []RoleResponse         `json:"roles"`
	Pager handlers.PagerResponse `json:"pager"`
}
