package attachmentservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
	"github.com/aportela/doneo/internal/repositories/historyoperationrepository"
	"github.com/aportela/doneo/internal/utils"
)

type AttachmentService interface {
	GetAttachment(ctx context.Context, id string) (domain.Attachment, error)
	AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) (domain.Attachment, error)
	DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error)
}

type attachmentService struct {
	database   database.Database
	repository attachmentrepository.AttachmentRepository
}

func NewService(database database.Database, repository attachmentrepository.AttachmentRepository) AttachmentService {
	return &attachmentService{database: database, repository: repository}
}

func (service *attachmentService) GetAttachment(ctx context.Context, id string) (domain.Attachment, error) {
	attachment, err := service.repository.GetAttachment(ctx, id)
	if err != nil {
		return domain.Attachment{}, fmt.Errorf("[AttachmentService] failed to get attachment with ID %s: %w", id, err)
	}
	return attachment, nil
}
func (service *attachmentService) AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) (domain.Attachment, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Attachment{}, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Attachment{}, fmt.Errorf("[AttachmentService] user ID not found in context")
	}
	attachment.CreatedBy.ID = currentUserId
	attachment.CreatedAt = time.Now()
	err = attachmentrepository.NewRepository(service.database).AddAttachment(ctx, attachment)
	if err != nil {
		return domain.Attachment{}, err
	}
	err = service.repository.AddProjectAttachment(ctx, projectId, attachment.ID)
	if err != nil {
		return domain.Attachment{}, err
	}
	err = historyoperationrepository.NewRepository(service.database).AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: attachment.CreatedAt, OperationType: domain.EventProjectAttachmentAdded})
	if err != nil {
		return domain.Attachment{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Attachment{}, err
	}
	return attachment, nil
}

func (service *attachmentService) DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error {
	// TODO: remove data/attachments file from storage
	tx, err := service.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[AttachmentService] user ID not found in context")
	}
	err = service.repository.DeleteProjectAttachment(ctx, projectId, attachmentId)
	if err != nil {
		return err
	}
	err = service.repository.DeleteProjectAttachment(ctx, projectId, attachmentId)
	if err != nil {
		return err
	}
	err = historyoperationrepository.NewRepository(service.database).AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectAttachmentDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *attachmentService) GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error) {
	attachments, err := service.repository.GetProjectAttachments(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[AttachmentService] failed to get project attachments: %w", err)
	}
	return attachments, nil
}
