package attachmentservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
)

type AttachmentService interface {
	GetAttachment(ctx context.Context, id string) (domain.Attachment, error)
	AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error
	DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error)
}

type attachmentService struct {
	repository attachmentrepository.AttachmentRepository
}

func NewService(repository attachmentrepository.AttachmentRepository) AttachmentService {
	return &attachmentService{repository: repository}
}

func (service *attachmentService) GetAttachment(ctx context.Context, id string) (domain.Attachment, error) {
	attachment, err := service.repository.GetAttachment(ctx, id)
	if err != nil {
		return domain.Attachment{}, fmt.Errorf("[AttachmentService] failed to get attachment with ID %s: %w", id, err)
	}
	return attachmentrepository.DTOToDomain(attachment), nil
}
func (service *attachmentService) AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error {
	return service.repository.AddProjectAttachment(ctx, projectId, attachmentrepository.DomainToDTO(attachment))
}

func (service *attachmentService) DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error {
	// TODO: remove data/attachments
	return service.repository.DeleteProjectAttachment(ctx, projectId, attachmentId)
}

func (service *attachmentService) GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error) {
	attachments, err := service.repository.GetProjectAttachments(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[AttachmentService] failed to get project attachments: %w", err)
	}
	return attachmentrepository.DTOArrayToDomainArray(attachments), nil
}
