package profilehandler

import (
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type updateRequest struct {
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password *string `json:"password,omitempty"`
}

type profileResponse struct {
	userhandler.UserBaseResponse
	Email     string `json:"email"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt *int64 `json:"updatedAt"`
}
