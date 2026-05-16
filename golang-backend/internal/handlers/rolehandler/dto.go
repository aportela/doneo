package rolehandler

import "github.com/aportela/doneo/internal/handlers"

type PermissionsFlags struct {
	AllowCreate  bool `json:"allowCreate"`
	AllowUpdate  bool `json:"allowUpdate"`
	AllowDelete  bool `json:"allowDelete"`
	AllowView    bool `json:"allowView"`
	AllowList    bool `json:"allowList"`
	AllowExecute bool `json:"allowExecute"`
}

type addRequest struct {
	Name        string           `json:"name"`
	Permissions PermissionsFlags `json:"permissions"`
}

type updateRequest struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions PermissionsFlags `json:"permissions"`
}

type SearchPermissionsFlags struct {
	AllowCreate  *bool `json:"allowCreate"`
	AllowUpdate  *bool `json:"allowUpdate"`
	AllowDelete  *bool `json:"allowDelete"`
	AllowView    *bool `json:"allowView"`
	AllowList    *bool `json:"allowList"`
	AllowExecute *bool `json:"allowExecute"`
}

type FilterRequest struct {
	Name        *string                 `json:"name"`
	Permissions *SearchPermissionsFlags `json:"permissions"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *FilterRequest        `json:"filter"`
}

type roleResponse struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Permissions PermissionsFlags `json:"permissions"`
}

type searchResponse struct {
	Roles []roleResponse         `json:"roles"`
	Pager handlers.PagerResponse `json:"pager"`
}
