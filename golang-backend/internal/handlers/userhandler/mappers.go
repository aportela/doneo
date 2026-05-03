package userhandler

import (
	"github.com/aportela/doneo/internal/domain"
)

func mapAddUserRequestToUserDomain(request addUserRequest) domain.User {
	return domain.User{
		UserBase:    domain.UserBase{ID: request.ID, Name: request.Name},
		Email:       request.Email,
		IsSuperUser: request.IsSuperUser,
	}
}

func mapUpdateUserRequestToUserDomain(request updateUserRequest) domain.User {
	return domain.User{
		UserBase:    domain.UserBase{ID: request.ID, Name: request.Name},
		Email:       request.Email,
		IsSuperUser: request.IsSuperUser,
	}
}

func mapUserDomainToUserResponse(user domain.User) userResponse {
	return userResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsSuperUser: user.IsSuperUser,
		AvatarURL:   user.AvatarURL,
	}
}

func mapUserDomainToAddUserResponse(user domain.User) addUserResponse {
	return addUserResponse{
		User: mapUserDomainToUserResponse(user),
	}
}

func mapUserDomainToUpdateUserResponse(user domain.User) updateUserResponse {
	return updateUserResponse{
		User: mapUserDomainToUserResponse(user),
	}
}

func mapUserDomainToGetUserResponse(user domain.User) getUserResponse {
	return getUserResponse{
		User: mapUserDomainToUserResponse(user),
	}
}

func mapUserArrayDomainToUserArrayResponse(users []domain.User) []userResponse {
	userResponses := []userResponse{}
	for _, user := range users {
		userResponses = append(userResponses, mapUserDomainToUserResponse(user))
	}
	return userResponses
}

func mapUserArrayDomainToSearchUsersResponse(users []domain.User) searchUsersResponse {
	return searchUsersResponse{
		Users: mapUserArrayDomainToUserArrayResponse(users),
	}
}
