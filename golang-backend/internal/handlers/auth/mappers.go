package handlers

import "github.com/aportela/doneo/internal/domain"

func ToUser(userRequest AuthRequest) domain.User {
	return domain.User{
		Email:    userRequest.Email,
		Password: &userRequest.Password,
	}
}
