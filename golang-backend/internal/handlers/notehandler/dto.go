package notehandler

import (
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type addRequest struct {
	Body string `json:"body"`
}

type updateRequest struct {
	ID   string `json:"id"`
	Body string `json:"body"`
}

type noteResponse struct {
	ID        string                       `json:"id"`
	CreatedBy userhandler.UserBaseResponse `json:"createdBy"`
	CreatedAt int64                        `json:"createdAt"`
	UpdatedAt *int64                       `json:"updatedAt"`
	Body      string                       `json:"body"`
}

type searchResponse struct {
	Notes []noteResponse `json:"notes"`
}
