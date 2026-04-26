package handlers

import "github.com/aportela/doneo/internal/domain"

func ToUser(userRequest UserRequest) domain.User {
	return domain.User{
		UserBase:    domain.UserBase{ID: userRequest.ID, Name: userRequest.Name},
		Email:       userRequest.Email,
		IsSuperUser: userRequest.IsSuperUser,
	}
}

func ToUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsSuperUser: user.IsSuperUser,
		AvatarURL:   user.AvatarURL,
	}
}

func ToUserResponses(users []domain.User) []UserResponse {
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, ToUserResponse(user))
	}
	return userResponses
}

func ToGetUserResponse(user domain.User) GetUserResponse {
	return GetUserResponse{User: ToUserResponse(user)}
}

func ToSearchUserResponse(users []domain.User) SearchUserResponse {
	return SearchUserResponse{Users: ToUserResponses(users)}
}
