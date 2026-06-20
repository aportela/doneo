package projectpermissionhandler

import (
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

// TODO: share userbase with user handler ??
type userBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TODO: share rolebase with role handler ??
type roleBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type addRequest struct {
	User userBase `json:"user"`
	Role roleBase `json:"role"`
}

type projectPermissionResponse struct {
	ID   string                       `json:"id"`
	User userhandler.UserBaseResponse `json:"user"`
	Role rolehandler.RoleResponse     `json:"role"`
}

type searchResponse struct {
	ProjectPermissions []projectPermissionResponse `json:"projectPermissions"`
}
