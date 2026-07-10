package pagehandler

import (
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type addRequest struct {
	Title string `json:"title"`
}

type updateRequest struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

type pageResponse struct {
	ID        string                       `json:"id"`
	CreatedBy userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt int64                        `json:"createdAt"`
	UpdatedAt *int64                       `json:"updatedAt"`
	Title     string                       `json:"title"`
	Body      string                       `json:"body"`
}

type searchResponse struct {
	Pages []pageResponse `json:"pages"`
}
