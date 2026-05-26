package rolehandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func requestPermissionsToDomainPermissionsBitmask(permissions permissionsFlags) domain.PermissionsBitmask {
	bitmaskPermission := domain.PermissionsBitmask(0)
	if permissions.AllowUpdateProject {
		bitmaskPermission.AddPermission(domain.PermissionUpdateProject)
	}
	if permissions.AllowDeleteProject {
		bitmaskPermission.AddPermission(domain.PermissionDeleteProject)
	}
	if permissions.AllowViewProject {
		bitmaskPermission.AddPermission(domain.PermissionViewProject)
	}
	if permissions.AllowAddTask {
		bitmaskPermission.AddPermission(domain.PermissionAddTask)
	}
	if permissions.AllowUpdateTask {
		bitmaskPermission.AddPermission(domain.PermissionUpdateTask)
	}
	if permissions.AllowDeleteTask {
		bitmaskPermission.AddPermission(domain.PermissionDeleteTask)
	}
	if permissions.AllowViewTask {
		bitmaskPermission.AddPermission(domain.PermissionViewTask)
	}
	return bitmaskPermission
}

func addRequestToDomain(request addRequest) domain.Role {
	return domain.Role{
		Name:               request.Name,
		PermissionsBitmask: requestPermissionsToDomainPermissionsBitmask(request.Permissions),
	}
}

func updateRequestToDomain(request updateRequest) domain.Role {
	return domain.Role{
		ID:                 request.Id,
		Name:               request.Name,
		PermissionsBitmask: requestPermissionsToDomainPermissionsBitmask(request.Permissions),
	}
}

func permissionDomainToResponsePermissionsFlags(bitmaskPermission domain.PermissionsBitmask) permissionsFlags {
	return permissionsFlags{
		AllowUpdateProject: bitmaskPermission.HasPermission(domain.PermissionUpdateProject),
		AllowDeleteProject: bitmaskPermission.HasPermission(domain.PermissionDeleteProject),
		AllowViewProject:   bitmaskPermission.HasPermission(domain.PermissionViewProject),
		AllowAddTask:       bitmaskPermission.HasPermission(domain.PermissionAddTask),
		AllowUpdateTask:    bitmaskPermission.HasPermission(domain.PermissionUpdateTask),
		AllowDeleteTask:    bitmaskPermission.HasPermission(domain.PermissionDeleteTask),
		AllowViewTask:      bitmaskPermission.HasPermission(domain.PermissionViewTask),
	}
}

func DomainToResponse(role domain.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissionDomainToResponsePermissionsFlags(role.PermissionsBitmask),
	}
}

func domainArrayToResponseArray(roles []domain.Role) []RoleResponse {
	roleResponses := []RoleResponse{}
	for _, role := range roles {
		roleResponses = append(roleResponses, DomainToResponse(role))
	}
	return roleResponses
}

func toSearchResponse(roles []domain.Role, pager browser.Result) searchResponse {
	return searchResponse{
		Roles: domainArrayToResponseArray(roles),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
