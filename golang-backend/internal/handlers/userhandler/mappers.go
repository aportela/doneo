package userhandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToUser(request addRequest) domain.User {
	return domain.User{
		UserBase:    domain.UserBase{Name: request.Name},
		Email:       request.Email,
		Password:    request.Password,
		IsSuperUser: request.IsSuperUser,
	}
}

func updateRequestToUser(request updateRequest) domain.User {
	return domain.User{
		UserBase:    domain.UserBase{Name: request.Name},
		Email:       request.Email,
		Password:    *request.Password,
		IsSuperUser: request.IsSuperUser,
	}
}

func userToResponse(user domain.User) userResponse {
	return userResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt.UnixMilli(),
		UpdatedAt:   utils.TimePtrToInt64Ptr(user.UpdatedAt),
		DeletedAt:   utils.TimePtrToInt64Ptr(user.DeletedAt),
		IsSuperUser: user.IsSuperUser,
		AvatarURL:   user.AvatarURL,
	}
}

func userToAddResponse(user domain.User) addResponse {
	return addResponse{
		User: userToResponse(user),
	}
}

func userToUpdateResponse(user domain.User) updateResponse {
	return updateResponse{
		User: userToResponse(user),
	}
}

func userToGetResponse(user domain.User) getResponse {
	return getResponse{
		User: userToResponse(user),
	}
}

func userArrayToResponse(users []domain.User) []userResponse {
	userResponses := []userResponse{}
	for _, user := range users {
		userResponses = append(userResponses, userToResponse(user))
	}
	return userResponses
}

func userArrayToSearchResponse(users []domain.User) searchResponse {
	return searchResponse{
		Users: userArrayToResponse(users),
	}
}
