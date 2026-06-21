package projectpermissionhandler

import (
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

type addRequest struct {
	User userhandler.UserBaseResponse `json:"user"`
	Role rolehandler.RoleBaseResponse `json:"role"`
}

type projectPermissionResponse struct {
	ID   string                       `json:"id"`
	User userhandler.UserBaseResponse `json:"user"`
	Role rolehandler.RoleBaseResponse `json:"role"`
}

type searchResponse struct {
	ProjectPermissions []projectPermissionResponse `json:"projectPermissions"`
}
