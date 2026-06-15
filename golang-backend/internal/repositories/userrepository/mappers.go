package userrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
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

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		if strings.Contains(sqlErr.Error(), "users.name") {
			return &domain.AlreadyExistsError{Field: "name"}
		} else if strings.Contains(sqlErr.Error(), "users.email") {
			return &domain.AlreadyExistsError{Field: "email"}
		}
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(name)") {
			return &domain.ValidationError{Field: "name"}
		} else if strings.Contains(sqlErr.Error(), "length(email)") {
			return &domain.ValidationError{Field: "email"}
		}
	}
	return err
}
