package userhandler

import "github.com/aportela/doneo/internal/handlers"

// TODO: this struct is duplicated on auth handler
type permissionsFlags struct {
	IsSuperUser bool `json:"isSuperUser"`
}

type addRequest struct {
	Name        string           `json:"name"`
	Email       string           `json:"email"`
	Password    string           `json:"password"`
	Permissions permissionsFlags `json:"permissions"`
}

type updateRequest struct {
	Id          string           `json:"id"`
	Name        string           `json:"name"`
	Email       string           `json:"email"`
	Password    *string          `json:"password,omitempty"`
	Permissions permissionsFlags `json:"permissions"`
}

type patchRequest struct {
	DeletedAt *int64 `json:"deletedAt"`
}

type filterUserPermissions struct {
	IsSuperUser *bool `json:"isSuperUser"`
}

type FilterRequest struct {
	Name        *string                   `json:"name"`
	Email       *string                   `json:"email"`
	Permissions *filterUserPermissions    `json:"permissions"`
	CreatedAt   *handlers.TimestampFilter `json:"createdAt"`
	UpdatedAt   *handlers.TimestampFilter `json:"updatedAt"`
	DeletedAt   *handlers.TimestampFilter `json:"deletedAt"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *FilterRequest        `json:"filter"`
}

type UserBaseResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type userResponse struct {
	UserBaseResponse
	Email       string           `json:"email"`
	CreatedAt   int64            `json:"createdAt"`
	UpdatedAt   *int64           `json:"updatedAt"`
	DeletedAt   *int64           `json:"deletedAt"`
	Permissions permissionsFlags `json:"permissions"`
}

type searchBaseResponse struct {
	Users []UserBaseResponse `json:"users"`
}

type searchResponse struct {
	Users []userResponse         `json:"users"`
	Pager handlers.PagerResponse `json:"pager"`
}
