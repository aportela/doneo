package attachmenthandler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/services/attachmentservice"
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
		if attachmentID, err := handler.service.SaveUploadedFile(file, header.Filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			attachment := domain.Attachment{
				ID:           attachmentID,
				OriginalName: header.Filename,
				ContentType:  header.Header.Get("Content-Type"),
				Size:         uint32(header.Size),
			}
			projectID := chi.URLParam(r, "project_id")
			if attachment, err := handler.service.AddProjectAttachment(r.Context(), projectID, attachment); err != nil {
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

func (handler *AttachmentHandler) DeleteProjectAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if attachment, err := handler.service.GetProjectAttachment(r.Context(), projectID, attachmentID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to get project attachment: %w", err))
		return
	} else {
		if err := handler.service.DeleteProjectAttachment(r.Context(), projectID, attachmentID); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete project attachment: %w", err))
			return
		}
		if err := handler.service.DeleteAttachment(attachment); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete project attachment storage file: %w", err))
			return
		}
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
}

func (handler *AttachmentHandler) DownloadProjectAttachment(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if attachment, err := handler.service.GetProjectAttachment(r.Context(), projectID, attachmentID); err != nil {
		// TODO: works with custom errors (like notFound / 404) ?
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		attachmentPath := handler.service.GetAttachmentFullPath(attachment.ID, attachment.OriginalName)
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
		if attachmentID, err := handler.service.SaveUploadedFile(file, header.Filename); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			attachment := domain.Attachment{
				ID:           attachmentID,
				OriginalName: header.Filename,
				ContentType:  header.Header.Get("Content-Type"),
				Size:         uint32(header.Size),
			}
			projectID := chi.URLParam(r, "project_id")
			taskID := chi.URLParam(r, "task_id")
			if attachment, err := handler.service.AddTaskAttachment(r.Context(), projectID, taskID, attachment); err != nil {
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

func (handler *AttachmentHandler) DeleteTaskAttachment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	attachmentID := chi.URLParam(r, "attachment_id")
	if attachment, err := handler.service.GetTaskAttachment(r.Context(), projectID, taskID, attachmentID); err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to get task attachment: %w", err))
		return
	} else {
		if err := handler.service.DeleteTaskAttachment(r.Context(), projectID, taskID, attachmentID); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete task attachment: %w", err))
			return
		}
		if err := handler.service.DeleteAttachment(attachment); err != nil {
			handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete task attachment storage file: %w", err))
			return
		}
		handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
	}
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
		attachmentPath := handler.service.GetAttachmentFullPath(attachment.ID, attachment.OriginalName)
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

func (handler *AttachmentHandler) GetTaskAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectID := chi.URLParam(r, "project_id")
	taskID := chi.URLParam(r, "task_id")
	projectAttachments, err := handler.service.GetTaskAttachments(r.Context(), projectID, taskID)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
