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
	Id          string           `json:"id"`
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

type roleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions permissionsFlags `json:"permissions"`
}

type searchResponse struct {
	Roles []roleResponse         `json:"roles"`
	Pager handlers.PagerResponse `json:"pager"`
}
