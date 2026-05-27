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
