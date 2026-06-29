package avatarhandler

import (
	"net/http"
	"strconv"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/avatarservice"
	"github.com/go-chi/chi/v5"
)

type AvatarHandler interface {
	GetAvatar(w http.ResponseWriter, r *http.Request)
	UploadAvatar(w http.ResponseWriter, r *http.Request)
	DeleteAvatar(w http.ResponseWriter, r *http.Request)
}

type avatarHandler struct {
	service                 avatarservice.AvatarService
	maxAvatarUploadFilesize int64
}

func NewHandler(service avatarservice.AvatarService, maxAvatarUploadFilesize int64) AvatarHandler {
	return &avatarHandler{service: service, maxAvatarUploadFilesize: maxAvatarUploadFilesize}
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

func (handler *avatarHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(handler.maxAvatarUploadFilesize); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if file, header, err := r.FormFile("file"); err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	} else {
		defer file.Close()
		if err := handler.service.SaveAvatar(r.Context(), file, header.Filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
}

func (handler *avatarHandler) DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	if err := handler.service.DeleteAvatar(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}
