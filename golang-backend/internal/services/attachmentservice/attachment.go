package attachmentservice

import (
	"context"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
)

type AttachmentService interface {
	AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error
	DeleteAttachment(ctx context.Context, attachmentId string) error
}

type attachmentService struct {
	repository attachmentrepository.AttachmentRepository
}

func NewAttachmentService(repository attachmentrepository.AttachmentRepository) AttachmentService {
	return &attachmentService{repository: repository}
}

func (s *attachmentService) AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error {
	return s.repository.AddProjectAttachment(ctx, projectId, attachmentrepository.DomainToDTO(attachment))
}

func (s *attachmentService) DeleteAttachment(ctx context.Context, attachmentId string) error {
	return s.repository.DeleteAttachment(ctx, attachmentId)
}
