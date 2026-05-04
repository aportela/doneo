package userrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func UserBaseToDTO(user domain.UserBase) userBaseDTO {
	return userBaseDTO{
		ID:   user.ID,
		Name: user.Name,
	}
}

func UserToDTO(user domain.User) userDTO {
	return userDTO{
		userBaseDTO:  UserBaseToDTO(user.UserBase),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.UnixMilli(),
		UpdatedAt:    utils.TimePtrToSQLNullInt64(user.UpdatedAt),
		DeletedAt:    utils.TimePtrToSQLNullInt64(user.DeletedAt),
		IsSuperUser:  user.IsSuperUser,
	}
}

func DTOToUserBase(user userBaseDTO) domain.UserBase {
	return domain.UserBase{
		ID:   user.ID,
		Name: user.Name,
	}
}

func DTOToUser(user userDTO) domain.User {
	return domain.User{
		UserBase:    DTOToUserBase(user.userBaseDTO),
		Email:       user.Email,
		CreatedAt:   time.UnixMilli(user.CreatedAt),
		UpdatedAt:   utils.SQLNullInt64ToTimePtr(user.UpdatedAt),
		DeletedAt:   utils.SQLNullInt64ToTimePtr(user.DeletedAt),
		IsSuperUser: user.IsSuperUser,
	}
}

func ToUserArray(users []userDTO) []domain.User {
	results := make([]domain.User, 0, len(users))
	for _, user := range users {
		results = append(results, DTOToUser(user))
	}
	return results
}
