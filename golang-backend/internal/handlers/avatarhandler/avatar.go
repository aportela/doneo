package avatarhandler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/avatarservice"
	"github.com/go-chi/chi/v5"
)

type AvatarHandler interface {
	Get(w http.ResponseWriter, r *http.Request)
	Add(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type avatarHandler struct {
	service                 avatarservice.AvatarService
	maxAvatarUploadFilesize int64
}

func NewHandler(service avatarservice.AvatarService, maxAvatarUploadFilesize int64) AvatarHandler {
	return &avatarHandler{service: service, maxAvatarUploadFilesize: maxAvatarUploadFilesize}
}

func (handler *avatarHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "user_id")
	if avatarPath, err := handler.service.GetUserAvatarPath(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		if info, err := os.Stat(avatarPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				w.Header().Set("Content-Type", "image/svg+xml")
				w.Write([]byte(avatarservice.DefaultAvatar))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			etag := fmt.Sprintf(`W/"%x-%x"`, info.ModTime().Unix(), info.Size())
			if r.Header.Get("If-None-Match") == etag {
				w.WriteHeader(http.StatusNotModified)
				return
			}
			w.Header().Set("ETag", etag)
			w.Header().Set("Content-Type", "image/svg+xml")
			http.ServeFile(w, r, avatarPath)
		}
	}
}

func (handler *avatarHandler) Add(w http.ResponseWriter, r *http.Request) {
	var request addRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[NoteHandler] invalid request payload: %w", err))
		return
	}
	if err := handler.service.SaveContextUserAvatar(r.Context(), request.Svg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *avatarHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := handler.service.DeleteContextUserAvatar(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}
