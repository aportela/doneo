package avatarhandler

import (
	"net/http"

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
	avatarSize := chi.URLParam(r, "avatar_size")
	http.Redirect(w, r, "https://i.pravatar.cc/"+avatarSize+"?u="+userID, http.StatusTemporaryRedirect)
}
