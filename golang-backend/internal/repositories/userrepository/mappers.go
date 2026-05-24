package userrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/utils"
)

func BaseDomainToBaseDTO(user domain.UserBase) UserBaseDTO {
	return UserBaseDTO{
		ID:   user.ID,
		Name: user.Name,
	}
}

func DomainToDTO(user domain.User) UserDTO {
	return UserDTO{
		UserBaseDTO:        BaseDomainToBaseDTO(user.UserBase),
		Email:              user.Email,
		PasswordHash:       user.PasswordHash,
		CreatedAt:          user.CreatedAt.UnixMilli(),
		UpdatedAt:          utils.TimePtrToSQLNullInt64(user.UpdatedAt),
		DeletedAt:          utils.TimePtrToSQLNullInt64(user.DeletedAt),
		PermissionsBitmask: uint64(user.PermissionsBitmask),
	}
}

func BaseDTOToBaseDomain(user UserBaseDTO) domain.UserBase {
	return domain.UserBase{
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: "https://i.pravatar.cc/32?u=" + user.ID,
	}
}

func DTOToDomain(user UserDTO) domain.User {
	return domain.User{
		UserBase:           BaseDTOToBaseDomain(user.UserBaseDTO),
		Email:              user.Email,
		CreatedAt:          time.UnixMilli(user.CreatedAt),
		UpdatedAt:          utils.SQLNullInt64ToTimePtr(user.UpdatedAt),
		DeletedAt:          utils.SQLNullInt64ToTimePtr(user.DeletedAt),
		PermissionsBitmask: domain.PermissionsBitmask(user.PermissionsBitmask),
	}
}

func DTOArrayToDomainArray(users []UserDTO) []domain.User {
	results := make([]domain.User, 0, len(users))
	for _, user := range users {
		results = append(results, DTOToDomain(user))
	}
	return results
}

func DomainFilterToDTO(filter domain.SearchUsersFilter) searchFilterDTO {
	return searchFilterDTO{
		Name:                        filter.Name,
		Email:                       filter.Email,
		RequiredPermissionsBitmask:  (*uint64)(filter.RequiredPermissionsBitmask),
		ForbiddenPermissionsBitmask: (*uint64)(filter.ForbiddenPermissionsBitmask),
		CreatedAt:                   repositories.TimestampFilterToDTO(filter.CreatedAt),
		UpdatedAt:                   repositories.TimestampFilterToDTO(filter.UpdatedAt),
		DeletedAt:                   repositories.TimestampFilterToDTO(filter.DeletedAt),
	}
}
