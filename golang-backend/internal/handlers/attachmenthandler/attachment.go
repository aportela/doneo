package attachmenthandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
		CreatedAt:    time.Now(),
		CreatedBy:    domain.UserBase{},
	}

	ext := filepath.Ext(header.Filename)

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

	projectId := chi.URLParam(r, "id")

	err = handler.service.AddProjectAttachment(r.Context(), projectId, attachment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	attachment, err = handler.service.GetAttachment(r.Context(), attachment.ID)
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
	projectId := chi.URLParam(r, "id")
	attachmentId := chi.URLParam(r, "attachment_id")
	err := handler.service.DeleteProjectAttachment(r.Context(), projectId, attachmentId)
	if err != nil {
		handlers.ToHandlerJSONResponse(w, nil, fmt.Errorf("[AttachmentHandler] failed to delete project attachment: %w", err))
		return
	}
	handlers.ToHandlerJSONResponse(w, handlers.ToEmptyResponse(), nil)
}

func (handler *AttachmentHandler) DownloadProjectAttachment(w http.ResponseWriter, r *http.Request) {

	//projectId := chi.URLParam(r, "id")
	attachmentId := chi.URLParam(r, "attachment_id")

	attachment, err := handler.service.GetAttachment(r.Context(), attachmentId)
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

func (handler *AttachmentHandler) GetProjectAttachments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	projectId := chi.URLParam(r, "id")
	projectAttachments, err := handler.service.GetProjectAttachments(r.Context(), projectId)
	handlers.ToHandlerJSONResponse(w, toSearchResponse(projectAttachments), err)
}
