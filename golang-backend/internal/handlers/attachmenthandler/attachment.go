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
	// 32 MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if file, header, err := r.FormFile("file"); err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	} else {
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
		if err := os.MkdirAll(dir, 0755); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fullPath := filepath.Join(dir, filename)
		if dst, err := os.Create(fullPath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			projectID := chi.URLParam(r, "project_id")
			if attachment, err := handler.service.AddProjectAttachment(r.Context(), projectID, attachment); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				if attachment, err := handler.service.GetProjectAttachment(r.Context(), projectID, attachment.ID); err != nil {
					// TODO: works with custom errors (like notFound / 404) ?
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					handlers.ToHandlerJSONResponse(w, domainToResponse(attachment), nil)
				}
			}
		}
	}
}

func (handler *AttachmentHandler) DeleteProjectAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if err := handler.service.DeleteProjectAttachment(r.Context(), projectID, attachmentID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete project attachment: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *AttachmentHandler) DownloadProjectAttachment(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if attachment, err := handler.service.GetProjectAttachment(r.Context(), projectID, attachmentID); err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		ext := filepath.Ext(attachment.OriginalName)
		filename := attachment.ID + ext
		dir := filepath.Join(
			handler.basePath,
			string(attachment.ID[len(attachment.ID)-2]),
			string(attachment.ID[len(attachment.ID)-1]),
		)
		attachmentPath := filepath.Join(dir, filename)
		if _, err = os.Stat(attachmentPath); err == nil {
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
}

func (handler *AttachmentHandler) GetProjectAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	projectAttachments, err := handler.service.GetProjectAttachments(r.Context(), projectID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}

func (handler *AttachmentHandler) AddTaskAttachment(w http.ResponseWriter, r *http.Request) {
	// 32 MB
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if file, header, err := r.FormFile("file"); err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	} else {
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
		if err = os.MkdirAll(dir, 0755); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fullPath := filepath.Join(dir, filename)
		if dst, err := os.Create(fullPath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			defer dst.Close()
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			projectID := chi.URLParam(r, "project_id")
			taskID := chi.URLParam(r, "task_id")
			if attachment, err := handler.service.AddTaskAttachment(r.Context(), projectID, taskID, attachment); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				if attachment, err := handler.service.GetTaskAttachment(r.Context(), projectID, taskID, attachment.ID); err != nil {
					fmt.Println(err)
					// TODO: works with custom errors (like notFound / 404) ?
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					handlers.ToHandlerJSONResponse(w, domainToResponse(attachment), nil)
				}
			}
		}
	}
}

func (handler *AttachmentHandler) DeleteTaskAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if err := handler.service.DeleteTaskAttachment(r.Context(), projectID, taskID, attachmentID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete task attachment: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *AttachmentHandler) DownloadTaskAttachment(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if attachment, err := handler.service.GetTaskAttachment(r.Context(), projectID, taskID, attachmentID); err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		ext := filepath.Ext(attachment.OriginalName)
		filename := attachment.ID + ext
		dir := filepath.Join(
			handler.basePath,
			string(attachment.ID[len(attachment.ID)-2]),
			string(attachment.ID[len(attachment.ID)-1]),
		)
		attachmentPath := filepath.Join(dir, filename)
		if _, err := os.Stat(attachmentPath); err == nil {
			// TODO: allow previews
			w.Header().Set("Content-Disposition", `attachment; filename="`+attachment.OriginalName+`"`)
			w.Header().Set("Content-Type", attachment.ContentType)
			http.ServeFile(w, r, attachmentPath)
		} else {
			// TODO: register error ?, attachment found in database but missing on disk
			handlers.ToHandlerJSONResponse(w, nil, err, http.StatusNotFound)
		}
	}
}

func (handler *AttachmentHandler) GetTaskAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	projectAttachments, err := handler.service.GetTaskAttachments(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
