package rolehandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func requestPermissionsToDomain(permissions PermissionFlags) domain.PermissionBitmask {
	var bitmaskPermission domain.PermissionBitmask
	bitmaskPermission = 0
	if permissions.AllowCreate {
		bitmaskPermission.AddPermission(domain.PermissionCreate)
	}
	if permissions.AllowUpdate {
		bitmaskPermission.AddPermission(domain.PermissionUpdate)
	}
	if permissions.AllowDelete {
		bitmaskPermission.AddPermission(domain.PermissionDelete)
	}
	if permissions.AllowView {
		bitmaskPermission.AddPermission(domain.PermissionView)
	}
	if permissions.AllowList {
		bitmaskPermission.AddPermission(domain.PermissionList)
	}
	if permissions.AllowExecute {
		bitmaskPermission.AddPermission(domain.PermissionExecute)
	}
	return bitmaskPermission
}

func addRequestToRole(request addRequest) domain.Role {
	return domain.Role{
		Name:       request.Name,
		Permission: requestPermissionsToDomain(request.Permissions),
	}
}

func updateRequestToRole(request updateRequest) domain.Role {
	return domain.Role{
		ID:         request.Id,
		Name:       request.Name,
		Permission: requestPermissionsToDomain(request.Permissions),
	}
}

func domainPermissionToResponsePermissions(bitmaskPermission domain.PermissionBitmask) PermissionFlags {
	return PermissionFlags{
		AllowCreate:  bitmaskPermission.HasPermission(domain.PermissionCreate),
		AllowUpdate:  bitmaskPermission.HasPermission(domain.PermissionUpdate),
		AllowDelete:  bitmaskPermission.HasPermission(domain.PermissionDelete),
		AllowView:    bitmaskPermission.HasPermission(domain.PermissionView),
		AllowList:    bitmaskPermission.HasPermission(domain.PermissionList),
		AllowExecute: bitmaskPermission.HasPermission(domain.PermissionExecute),
	}
}

func roleToResponse(role domain.Role) roleResponse {
	return roleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: domainPermissionToResponsePermissions(role.Permission),
	}
}

func roleArrayToResponse(roles []domain.Role) []roleResponse {
	roleResponses := []roleResponse{}
	for _, role := range roles {
		roleResponses = append(roleResponses, roleToResponse(role))
	}
	return roleResponses
}

func ToSearchResponse(roles []domain.Role, pager browser.Result) searchResponse {
	return searchResponse{
		Roles: roleArrayToResponse(roles),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
