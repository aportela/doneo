package userrepository

import (
	"github.com/aportela/doneo/internal/domain"
)

func MapUserDomainToUserDTO(user domain.User) userDTO {
	return userDTO{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		DeletedAt:    user.DeletedAt,
		IsSuperUser:  user.IsSuperUser,
	}
}

func MapUserDTOToUserDomain(user userDTO) domain.User {
	return domain.User{
		UserBase: domain.UserBase{
			ID:   user.ID,
			Name: user.Name,
		},
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
		IsSuperUser: user.IsSuperUser,
	}
}

func MapUserArrayDTOToUserArrayDomain(users []userDTO) []domain.User {
	results := []domain.User{}
	for _, user := range users {
		results = append(results, MapUserDTOToUserDomain(user))
	}
	return results
}
