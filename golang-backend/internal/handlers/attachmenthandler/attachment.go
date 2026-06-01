package attachmenthandler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
	"github.com/aportela/doneo/internal/services/attachmentservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/go-chi/chi/v5"
)

type AttachmentHandler struct {
	service  attachmentservice.AttachmentService
	basePath string
}

func NewAttachmentHandler(db database.Database, basePath string) *AttachmentHandler {
	repository := attachmentrepository.NewAttachmentRepository(db)
	service := attachmentservice.NewAttachmentService(repository)
	return &AttachmentHandler{service: service, basePath: basePath}
}

func (h *AttachmentHandler) AddProjectAttachment(w http.ResponseWriter, r *http.Request) {
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
	attachment.CreatedBy.ID, _ = middlewares.GetUserIDFromContext(r.Context())

	ext := filepath.Ext(header.Filename)

	filename := attachment.ID + ext
	dir := filepath.Join(
		h.basePath,
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

	err = h.service.AddProjectAttachment(r.Context(), projectId, attachment)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	handlers.ToHandlerJSONResponse(w, domainToResponse(attachment), nil)
}
