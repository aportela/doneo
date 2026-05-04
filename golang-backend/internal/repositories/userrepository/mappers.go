package userrepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

func UserBaseToDTO(user domain.UserBase) UserBaseDTO {
	return UserBaseDTO{
		ID:   user.ID,
		Name: user.Name,
	}
}

func UserToDTO(user domain.User) UserDTO {
	return UserDTO{
		UserBaseDTO:  UserBaseToDTO(user.UserBase),
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    user.CreatedAt.UnixMilli(),
		UpdatedAt:    utils.TimePtrToSQLNullInt64(user.UpdatedAt),
		DeletedAt:    utils.TimePtrToSQLNullInt64(user.DeletedAt),
		IsSuperUser:  user.IsSuperUser,
	}
}

func DTOToUserBase(user UserBaseDTO) domain.UserBase {
	return domain.UserBase{
		ID:        user.ID,
		Name:      user.Name,
		AvatarURL: "https://i.pravatar.cc/32?u=" + user.ID,
	}
}

func DTOToUser(user UserDTO) domain.User {
	return domain.User{
		UserBase:    DTOToUserBase(user.UserBaseDTO),
		Email:       user.Email,
		CreatedAt:   time.UnixMilli(user.CreatedAt),
		UpdatedAt:   utils.SQLNullInt64ToTimePtr(user.UpdatedAt),
		DeletedAt:   utils.SQLNullInt64ToTimePtr(user.DeletedAt),
		IsSuperUser: user.IsSuperUser,
	}
}

func ToUserArray(users []UserDTO) []domain.User {
	results := make([]domain.User, 0, len(users))
	for _, user := range users {
		results = append(results, DTOToUser(user))
	}
	return results
}
