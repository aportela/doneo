package projectpermissionhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.ProjectPermission {
	return domain.ProjectPermission{
		User: domain.UserBase{
			ID: request.UserId,
		},
		Role: domain.Role{
			ID: request.RoleId,
		},
	}
}

func domainToResponse(projectPermission domain.ProjectPermission) ProjectPermissionResponse {
	return ProjectPermissionResponse{
		ID:   projectPermission.ID,
		User: userhandler.BaseDomainToBaseResponse(projectPermission.User),
		Role: rolehandler.DomainToResponse(projectPermission.Role),
	}
}

func domainArrayToResponseArray(projectPermissions []domain.ProjectPermission) []ProjectPermissionResponse {
	projectPermissionsResponse := []ProjectPermissionResponse{}
	for _, projectPermission := range projectPermissions {
		projectPermissionsResponse = append(projectPermissionsResponse, domainToResponse(projectPermission))
	}
	return projectPermissionsResponse
}

func toSearchResponse(projectPermissions []domain.ProjectPermission) searchResponse {
	return searchResponse{
		ProjectPermissions: domainArrayToResponseArray(projectPermissions),
	}
}
