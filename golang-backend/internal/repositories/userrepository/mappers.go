package userrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/utils"
)

func toBaseDTO(user domain.UserBase) userBaseDTO {
	return userBaseDTO{
		ID:   user.ID,
		Name: user.Name,
	}
}

func toDTO(user domain.User) userDTO {
	return userDTO{
		userBaseDTO:        toBaseDTO(user.UserBase),
		Email:              user.Email,
		CreatedAt:          user.CreatedAt.UnixMilli(),
		UpdatedAt:          utils.TimePtrToSQLNullInt64(user.UpdatedAt),
		DeletedAt:          utils.TimePtrToSQLNullInt64(user.DeletedAt),
		PermissionsBitmask: uint64(user.PermissionsBitmask),
	}
}

func toBaseDomain(user userBaseDTO) domain.UserBase {
	return domain.UserBase{
		ID:   user.ID,
		Name: user.Name,
	}
}

func toDomain(user userDTO) domain.User {
	return domain.User{
		UserBase:           toBaseDomain(user.userBaseDTO),
		Email:              user.Email,
		PasswordHash:       user.PasswordHash,
		CreatedAt:          time.UnixMilli(user.CreatedAt),
		UpdatedAt:          utils.SQLNullInt64ToTimePtr(user.UpdatedAt),
		DeletedAt:          utils.SQLNullInt64ToTimePtr(user.DeletedAt),
		PermissionsBitmask: domain.Bitmask(user.PermissionsBitmask),
	}
}

func toDomainArray(users []userDTO) []domain.User {
	results := make([]domain.User, 0, len(users))
	for _, user := range users {
		results = append(results, toDomain(user))
	}
	return results
}

func toFilterDTO(filter domain.SearchUsersFilter) searchFilterDTO {
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
