package projectpermissionhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
)

func addRequestToDomain(request addRequest) domain.ProjectPermission {
	return domain.ProjectPermission{
		User: domain.UserBase{
			ID: request.User.ID,
		},
		Role: domain.Role{
			RoleBase: domain.RoleBase{
				ID: request.Role.ID,
			},
		},
	}
}

func domainToResponse(projectPermission domain.ProjectPermission) projectPermissionResponse {
	return projectPermissionResponse{
		ID:   projectPermission.ID,
		User: userhandler.BaseDomainToBaseResponse(projectPermission.User),
		Role: rolehandler.DomainToResponse(projectPermission.Role),
	}
}

func domainArrayToResponseArray(projectPermissions []domain.ProjectPermission) []projectPermissionResponse {
	projectPermissionsResponse := []projectPermissionResponse{}
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
