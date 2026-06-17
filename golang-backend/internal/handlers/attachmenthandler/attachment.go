package attachmenthandler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/attachmentservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type AttachmentHandler struct {
	service  attachmentservice.AttachmentService
	basePath string
}

func NewHandler(service attachmentservice.AttachmentService, basePath string) *AttachmentHandler {
	return &AttachmentHandler{service: service, basePath: basePath}
}

func (handler *AttachmentHandler) AddProjectAttachment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	attachment := domain.Attachment{
		ID:           utils.UUID(),
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
		Size:         uint32(header.Size),
	}

	ext := filepath.Ext(header.Filename)

	// TODO: move to service ???
	filename := attachment.ID + ext
	dir := filepath.Join(
		handler.basePath,
		string(attachment.ID[len(attachment.ID)-2]),
		string(attachment.ID[len(attachment.ID)-1]),
	)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullPath := filepath.Join(dir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectID := chi.URLParam(r, "project_id")

	attachment, err = handler.service.AddProjectAttachment(r.Context(), projectID, attachment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attachment, err = handler.service.GetProjectAttachment(r.Context(), projectID, attachment.ID)
	if err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	handlers.ToHandlerJSONResponse(w, domainToResponse(attachment), nil)
}

func (handler *AttachmentHandler) DeleteProjectAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	err := handler.service.DeleteProjectAttachment(r.Context(), projectID, attachmentID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete project attachment: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *AttachmentHandler) DownloadProjectAttachment(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")

	attachment, err := handler.service.GetProjectAttachment(r.Context(), projectID, attachmentID)
	if err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(attachment.OriginalName)

	filename := attachment.ID + ext

	dir := filepath.Join(
		handler.basePath,
		string(attachment.ID[len(attachment.ID)-2]),
		string(attachment.ID[len(attachment.ID)-1]),
	)

	attachmentPath := filepath.Join(dir, filename)

	_, err = os.Stat(attachmentPath)

	if err == nil {

		switch chi.URLParam(r, "mode") {
		case "download":
			w.Header().Set("Content-Disposition", `attachment; filename="`+url.PathEscape(attachment.OriginalName)+`"`)

		case "inline":
			w.Header().Set("Content-Disposition", `inline; filename="`+url.PathEscape(attachment.OriginalName)+`"`)
		default:
			// TODO:
		}
		w.Header().Set("Content-Type", attachment.ContentType)

		http.ServeFile(w, r, attachmentPath)
	} else {
		// TODO: register error ?, attachment found in database but missing on disk
		handlers.ToHandlerJSONResponse(w, nil, err, http.StatusNotFound)
	}
}

func (handler *AttachmentHandler) GetProjectAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	projectAttachments, err := handler.service.GetProjectAttachments(r.Context(), projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}

func (handler *AttachmentHandler) AddTaskAttachment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32 MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	attachment := domain.Attachment{
		ID:           utils.UUID(),
		OriginalName: header.Filename,
		ContentType:  header.Header.Get("Content-Type"),
		Size:         uint32(header.Size),
	}

	ext := filepath.Ext(header.Filename)

	// TODO: move to service ???
	filename := attachment.ID + ext
	dir := filepath.Join(
		handler.basePath,
		string(attachment.ID[len(attachment.ID)-2]),
		string(attachment.ID[len(attachment.ID)-1]),
	)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fullPath := filepath.Join(dir, filename)

	dst, err := os.Create(fullPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")

	attachment, err = handler.service.AddTaskAttachment(r.Context(), projectID, taskID, attachment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attachment, err = handler.service.GetProjectAttachment(r.Context(), projectID, attachment.ID)
	if err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	handlers.ToHandlerJSONResponse(w, domainToResponse(attachment), nil)
}

func (handler *AttachmentHandler) DeleteTaskAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	err := handler.service.DeleteTaskAttachment(r.Context(), projectID, taskID, attachmentID)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete task attachment: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *AttachmentHandler) DownloadTaskAttachment(w http.ResponseWriter, r *http.Request) {

	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	attachmentID := chi.URLParam(r, "attachment_id")

	attachment, err := handler.service.GetTaskAttachment(r.Context(), projectID, taskID, attachmentID)
	if err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(attachment.OriginalName)

	filename := attachment.ID + ext

	dir := filepath.Join(
		handler.basePath,
		string(attachment.ID[len(attachment.ID)-2]),
		string(attachment.ID[len(attachment.ID)-1]),
	)

	attachmentPath := filepath.Join(dir, filename)

	_, err = os.Stat(attachmentPath)

	if err == nil {

		// TODO: allow previews
		w.Header().Set("Content-Disposition", `attachment; filename="`+attachment.OriginalName+`"`)
		w.Header().Set("Content-Type", attachment.ContentType)

		http.ServeFile(w, r, attachmentPath)
	} else {
		// TODO: register error ?, attachment found in database but missing on disk
		handlers.ToHandlerJSONResponse(w, nil, err, http.StatusNotFound)
	}
}

func (handler *AttachmentHandler) GetTaskAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	projectAttachments, err := handler.service.GetTaskAttachments(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
