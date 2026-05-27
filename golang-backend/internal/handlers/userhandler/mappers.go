package userhandler

import (
	"github.com/aportela/doneo/internal/browser"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/utils"
)

func permissionsToDomainPermissionsBitmask(permissions permissionsFlags) domain.PermissionsBitmask {
	var permissionsBitmask domain.PermissionsBitmask
	if permissions.IsSuperUser {
		permissionsBitmask.AddPermission(domain.UserPermissionAdmin)
	}
	return permissionsBitmask
}

func addRequestToDomain(request addRequest) domain.User {
	return domain.User{
		UserBase:           domain.UserBase{Name: request.Name},
		Email:              request.Email,
		Password:           request.Password,
		PermissionsBitmask: permissionsToDomainPermissionsBitmask(request.Permissions),
	}
}

func updateRequestToDomain(request updateRequest) domain.User {
	user := domain.User{
		UserBase:           domain.UserBase{Name: request.Name},
		Email:              request.Email,
		PermissionsBitmask: permissionsToDomainPermissionsBitmask(request.Permissions),
	}
	if request.Password != nil {
		user.Password = *request.Password
	}
	return user
}

func permissionsDomainToResponsePermissionsFlags(permissionsBitmask domain.PermissionsBitmask) permissionsFlags {
	return permissionsFlags{
		IsSuperUser: permissionsBitmask.HasPermission(domain.UserPermissionAdmin),
	}
}

func domainToResponse(user domain.User) userResponse {
	return userResponse{
		UserBaseResponse: UserBaseResponse{
			ID:   user.ID,
			Name: user.Name,
		},
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.UnixMilli(),
		UpdatedAt:   utils.TimePtrToInt64Ptr(user.UpdatedAt),
		DeletedAt:   utils.TimePtrToInt64Ptr(user.DeletedAt),
		Permissions: permissionsDomainToResponsePermissionsFlags(user.PermissionsBitmask),
	}
}

func BaseDomainToBaseResponse(user domain.UserBase) UserBaseResponse {
	return UserBaseResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func domainToBaseResponse(user domain.User) UserBaseResponse {
	return UserBaseResponse{
		ID:   user.ID,
		Name: user.Name,
	}
}

func domainArrayToResponseArray(users []domain.User) []userResponse {
	userResponses := []userResponse{}
	for _, user := range users {
		userResponses = append(userResponses, domainToResponse(user))
	}
	return userResponses
}

func domainArrayToBaseResponseArray(users []domain.User) []UserBaseResponse {
	userResponses := []UserBaseResponse{}
	for _, user := range users {
		userResponses = append(userResponses, domainToBaseResponse(user))
	}
	return userResponses
}

func toSearchBaseResponse(users []domain.User) searchBaseResponse {
	return searchBaseResponse{
		Users: domainArrayToBaseResponseArray(users),
	}
}

func toSearchResponse(users []domain.User, pager browser.Result) searchResponse {
	return searchResponse{
		Users: domainArrayToResponseArray(users),
		Pager: handlers.PagerResponse{
			Enabled:      pager.ResultsPage > 0,
			CurrentPage:  pager.CurrentPage,
			ResultsPage:  pager.ResultsPage,
			TotalPages:   pager.TotalPages,
			TotalResults: pager.TotalResults,
		},
	}
}
