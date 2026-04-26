package handlers

import (
	"database/sql"
	"net/http"

	"github.com/aportela/doneo/internal/models"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	user := models.User{}
	err := h.service.AddUser(ctx, user)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	user := models.User{}
	err := h.service.UpdateUser(ctx, user)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ToJSONResponse(w, http.StatusNotFound, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		} else {
			utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, user)

}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	err := h.service.DeleteUser(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ToJSONResponse(w, http.StatusNotFound, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		} else {
			utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, nil)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	userId := chi.URLParam(r, "id")
	user, err := h.service.GetUser(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ToJSONResponse(w, http.StatusNotFound, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		} else {
			utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
				"debugErrorMessage": err.Error(),
			})
			return
		}
	}
	utils.ToJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")
	users, err := h.service.SearchUsers(ctx)
	if err != nil {
		utils.ToJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"debugErrorMessage": err.Error(),
		})
		return
	}
	utils.ToJSONResponse(w, http.StatusOK, users)
}
