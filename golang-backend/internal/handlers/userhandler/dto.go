package userhandler

import "github.com/aportela/doneo/internal/handlers"

type addRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsSuperUser bool   `json:"isSuperUser"`
}

type updateRequest struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Password    *string `json:"password,omitempty"`
	IsSuperUser bool    `json:"isSuperUser"`
}

type patchRequest struct {
	DeletedAt *int64 `json:"deletedAt"`
}

type TimestampFilter struct {
	From *int `json:"from"`
	To   *int `json:"to"`
}

type FilterRequest struct {
	Name              *string          `json:"name"`
	Email             *string          `json:"email"`
	AdministratorFlag *bool            `json:"administratorFlag"`
	CreatedAt         *TimestampFilter `json:"createdAt"`
	UpdatedAt         *TimestampFilter `json:"updatedAt"`
	DeletedAt         *TimestampFilter `json:"deletedAt"`
}

type searchRequest struct {
	Pager  handlers.PagerRequest `json:"pager"`
	Order  handlers.OrderRequest `json:"order"`
	Filter *FilterRequest        `json:"filter"`
}

type userResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	CreatedAt   int64  `json:"createdAt"`
	UpdatedAt   *int64 `json:"updatedAt"`
	DeletedAt   *int64 `json:"deletedAt"`
	IsSuperUser bool   `json:"isSuperUser"`
	AvatarURL   string `json:"avatarUrl"`
}

type searchResponse struct {
	Users []userResponse         `json:"users"`
	Pager handlers.PagerResponse `json:"pager"`
}
