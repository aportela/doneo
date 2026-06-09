package rolehandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
)

func requestPermissionsToDomainPermissionsBitmask(permissions permissionsFlags) domain.Bitmask {
	bitmaskPermission := domain.Bitmask(0)
	if permissions.AllowUpdateProject {
		bitmaskPermission.AddFlag(domain.PermissionUpdateProject)
	}
	if permissions.AllowDeleteProject {
		bitmaskPermission.AddFlag(domain.PermissionDeleteProject)
	}
	if permissions.AllowViewProject {
		bitmaskPermission.AddFlag(domain.PermissionViewProject)
	}
	if permissions.AllowAddTask {
		bitmaskPermission.AddFlag(domain.PermissionAddTask)
	}
	if permissions.AllowUpdateTask {
		bitmaskPermission.AddFlag(domain.PermissionUpdateTask)
	}
	if permissions.AllowDeleteTask {
		bitmaskPermission.AddFlag(domain.PermissionDeleteTask)
	}
	if permissions.AllowViewTask {
		bitmaskPermission.AddFlag(domain.PermissionViewTask)
	}
	return bitmaskPermission
}

func addRequestToDomain(request addRequest) domain.Role {
	return domain.Role{
		RoleBase: domain.RoleBase{
			Name: request.Name,
		},
		PermissionsBitmask: requestPermissionsToDomainPermissionsBitmask(request.Permissions),
	}
}

func updateRequestToDomain(request updateRequest) domain.Role {
	return domain.Role{
		RoleBase: domain.RoleBase{
			ID:   request.Id,
			Name: request.Name,
		},
		PermissionsBitmask: requestPermissionsToDomainPermissionsBitmask(request.Permissions),
	}
}

func permissionDomainToResponsePermissionsFlags(bitmaskPermission domain.Bitmask) permissionsFlags {
	return permissionsFlags{
		AllowUpdateProject: bitmaskPermission.HasFlag(domain.PermissionUpdateProject),
		AllowDeleteProject: bitmaskPermission.HasFlag(domain.PermissionDeleteProject),
		AllowViewProject:   bitmaskPermission.HasFlag(domain.PermissionViewProject),
		AllowAddTask:       bitmaskPermission.HasFlag(domain.PermissionAddTask),
		AllowUpdateTask:    bitmaskPermission.HasFlag(domain.PermissionUpdateTask),
		AllowDeleteTask:    bitmaskPermission.HasFlag(domain.PermissionDeleteTask),
		AllowViewTask:      bitmaskPermission.HasFlag(domain.PermissionViewTask),
	}
}

func DomainToResponse(role domain.Role) RoleResponse {
	return RoleResponse{
		RoleBaseResponse: RoleBaseResponse{
			ID:   role.ID,
			Name: role.Name,
		},
		Permissions: permissionDomainToResponsePermissionsFlags(role.PermissionsBitmask),
	}
}

func baseDomainToBaseResponse(role domain.RoleBase) RoleBaseResponse {
	return RoleBaseResponse{
		ID:   role.ID,
		Name: role.Name,
	}
}

func domainToBaseResponse(role domain.Role) RoleBaseResponse {
	return RoleBaseResponse{
		ID:   role.ID,
		Name: role.Name,
	}
}

func domainArrayToResponseArray(roles []domain.Role) []RoleResponse {
	roleResponses := []RoleResponse{}
	for _, role := range roles {
		roleResponses = append(roleResponses, DomainToResponse(role))
	}
	return roleResponses
}

func domainArrayToBaseResponseArray(roles []domain.Role) []RoleBaseResponse {
	roleResponses := []RoleBaseResponse{}
	for _, role := range roles {
		roleResponses = append(roleResponses, domainToBaseResponse(role))
	}
	return roleResponses
}

func toSearchBaseResponse(roles []domain.Role) searchBaseResponse {
	return searchBaseResponse{
		Roles: domainArrayToBaseResponseArray(roles),
	}
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
