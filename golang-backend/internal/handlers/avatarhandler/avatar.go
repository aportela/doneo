package avatarhandler

import (
	"net/http"
	"strconv"

	"github.com/aportela/doneo/internal/services/avatarservice"
	"github.com/go-chi/chi/v5"
)

type AvatarHandler interface {
	GetAvatar(w http.ResponseWriter, r *http.Request)
}

type avatarHandler struct {
	service avatarservice.AvatarService
}

func NewHandler(service avatarservice.AvatarService) AvatarHandler {
	return &avatarHandler{service: service}
}

func (handler *avatarHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	var avatarSize avatarservice.AvatarSize
	switch chi.URLParam(r, "avatar_size") {
	case "tiny":
		avatarSize = avatarservice.AvatarSizeTiny
	case "small":
		avatarSize = avatarservice.AvatarSizeSmall
	case "normal":
		avatarSize = avatarservice.AvatarSizeNormal
	}

	http.Redirect(w, r, "https://i.pravatar.cc/"+strconv.FormatUint(uint64(avatarSize), 10)+"?u="+userID, http.StatusTemporaryRedirect)
}
