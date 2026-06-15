package attachmentservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/attachmentrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type AttachmentService interface {
	GetAttachment(ctx context.Context, attachmentID string) (domain.Attachment, error)

	AddProjectAttachment(ctx context.Context, projectID string, attachment domain.Attachment) (domain.Attachment, error)
	DeleteProjectAttachment(ctx context.Context, projectID string, attachmentID string) error
	GetProjectAttachments(ctx context.Context, projectID string) ([]domain.Attachment, error)

	AddTaskAttachment(ctx context.Context, projectID string, taskID string, attachment domain.Attachment) (domain.Attachment, error)
	DeleteTaskAttachment(ctx context.Context, projectID string, taskID string, attachmentID string) error
	GetTaskAttachments(ctx context.Context, projectID string, taskID string) ([]domain.Attachment, error)
}

type attachmentService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	attachmentRepository    attachmentrepository.AttachmentRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, attachmentRepository attachmentrepository.AttachmentRepository) AttachmentService {
	return &attachmentService{db: db, authorizationService: authorizationService, historyOperationService: historyOperationService, attachmentRepository: attachmentRepository}
}

func (service *attachmentService) withProjectUpdatePermission(ctx context.Context, projectID string, action func(currentUserID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	if err := service.authorizationService.RequireProjectUpdatePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *attachmentService) GetAttachment(ctx context.Context, attachmentID string) (domain.Attachment, error) {
	// TODO: check permissions
	attachment, err := service.attachmentRepository.GetAttachment(ctx, service.db, attachmentID)
	if err != nil {
		return domain.Attachment{}, fmt.Errorf("[AttachmentService] failed to get attachment with ID %s: %w", attachmentID, err)
	}
	return attachment, nil
}

func (service *attachmentService) AddProjectAttachment(ctx context.Context, projectID string, attachment domain.Attachment) (domain.Attachment, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		attachment.CreatedBy.ID = currentUserID
		attachment.CreatedAt = time.Now()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.attachmentRepository.AddAttachment(ctx, tx, attachment); err != nil {
				return err
			}
			if err := service.attachmentRepository.AddProjectAttachment(ctx, tx, projectID, attachment.ID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     attachment.CreatedAt,
					OperationType: domain.EventProjectAttachmentAdded,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		return domain.Attachment{}, err
	}
	return attachment, nil
}

func (service *attachmentService) DeleteProjectAttachment(ctx context.Context, projectID string, attachmentID string) error {
	// TODO: remove data/attachments file from storage
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.attachmentRepository.DeleteProjectAttachment(ctx, tx, projectID, attachmentID); err != nil {
				return err
			}
			if err := service.attachmentRepository.DeleteAttachment(ctx, tx, attachmentID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectAttachmentDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	return err
}

func (service *attachmentService) GetProjectAttachments(ctx context.Context, projectID string) ([]domain.Attachment, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	attachments, err := service.attachmentRepository.GetProjectAttachments(ctx, service.db, projectID)
	if err != nil {
		return nil, fmt.Errorf("[AttachmentService] failed to get project attachments: %w", err)
	}
	return attachments, nil
}

func (service *attachmentService) AddTaskAttachment(ctx context.Context, projectID string, taskID string, attachment domain.Attachment) (domain.Attachment, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		attachment.CreatedBy.ID = currentUserID
		attachment.CreatedAt = time.Now()
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.attachmentRepository.AddAttachment(ctx, tx, attachment); err != nil {
				return err
			}
			if err := service.attachmentRepository.AddTaskAttachment(ctx, tx, taskID, attachment.ID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     attachment.CreatedAt,
					OperationType: domain.EventTaskAttachmentAdded,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	if err != nil {
		return domain.Attachment{}, err
	}
	return attachment, nil
}

func (service *attachmentService) DeleteTaskAttachment(ctx context.Context, projectID string, taskID string, attachmentID string) error {
	// TODO: remove data/attachments file from storage
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.attachmentRepository.DeleteTaskAttachment(ctx, tx, taskID, attachmentID); err != nil {
				return err
			}
			if err := service.attachmentRepository.DeleteAttachment(ctx, tx, attachmentID); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskAttachmentDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	})
	return err
}

func (service *attachmentService) GetTaskAttachments(ctx context.Context, projectID string, taskID string) ([]domain.Attachment, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	attachments, err := service.attachmentRepository.GetTaskAttachments(ctx, service.db, taskID)
	if err != nil {
		return nil, fmt.Errorf("[AttachmentService] failed to get task attachments: %w", err)
	}
	return attachments, nil
}
