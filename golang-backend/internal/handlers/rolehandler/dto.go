package rolehandler

import "github.com/aportela/doneo/internal/handlers"

type PermissionFlags struct {
	AllowCreate  bool `json:"allowCreate"`
	AllowUpdate  bool `json:"allowUpdate"`
	AllowDelete  bool `json:"allowDelete"`
	AllowView    bool `json:"allowView"`
	AllowList    bool `json:"allowList"`
	AllowExecute bool `json:"allowExecute"`
}

type addRequest struct {
	Name        string          `json:"name"`
	Permissions PermissionFlags `json:"permissions"`
}

type updateRequest struct {
	Id          string          `json:"id"`
	Name        string          `json:"name"`
	Permissions PermissionFlags `json:"permissions"`
}

type searchRequest struct {
	Pager handlers.PagerRequest `json:"pager"`
	Order handlers.OrderRequest `json:"order"`
}

type roleResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Permissions PermissionFlags `json:"permissions"`
}

type searchResponse struct {
	Roles []roleResponse         `json:"roles"`
	Pager handlers.PagerResponse `json:"pager"`
}
