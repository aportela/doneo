package profilehandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func updateRequestToDomain(request updateRequest) domain.User {
	user := domain.User{
		UserBase: domain.UserBase{Name: request.Name},
		Email:    request.Email,
	}
	return user
}

func domainToResponse(user domain.User) profileResponse {
	return profileResponse{
		UserBaseResponse: userhandler.UserBaseResponse{
			ID:   user.ID,
			Name: user.Name,
		},
		Email:     user.Email,
		CreatedAt: user.CreatedAt.UnixMilli(),
		UpdatedAt: utils.TimePtrToInt64Ptr(user.UpdatedAt),
	}
}
