package attachmentrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type AttachmentRepository interface {
	AddAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachment domain.Attachment) error
	DeleteAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachmentID string) error
	GetAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachmentID string) (domain.Attachment, error)

	AddProjectAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, attachmentID string) error
	DeleteProjectAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, attachmentID string) error
	GetProjectAttachments(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Attachment, error)

	AddTaskAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, attachmentID string) error
	DeleteTaskAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, attachmentID string) error
	GetTaskAttachments(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Attachment, error)
}

type attachmentRepository struct {
}

func NewRepository() AttachmentRepository {
	return &attachmentRepository{}
}

func (repository *attachmentRepository) AddAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachment domain.Attachment) error {
	dto := toDTO(attachment)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO attachments
				(id, original_name, content_type, size, creator_id, created_at)
			VALUES
				(?, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		dto.OriginalName,
		dto.ContentType,
		dto.Size,
		dto.CreatorID,
		dto.CreatedAt,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *attachmentRepository) DeleteAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachmentID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM attachments
			WHERE
				id = ?
        `,
		attachmentID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *attachmentRepository) GetAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, attachmentID string) (domain.Attachment, error) {
	var dto attachmentDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
				A.id, A.creator_id, U.name, A.created_at, A.original_name, A.content_type, A.size
            FROM project_attachments PA
			INNER JOIN attachments A ON A.id = PA.attachment_id
			INNER JOIN users U ON U.id = A.creator_id
            WHERE
				A.id = ?
        `,
		attachmentID).Scan(&dto.ID, &dto.CreatorID, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Attachment{}, domain.NotFoundError
		}
		return domain.Attachment{}, err
	}
	return toDomain(dto), err
}

func (repository *attachmentRepository) AddProjectAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, attachmentID string) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_attachments (project_id, attachment_id)
			VALUES (?, ?)
        `,
		projectID,
		attachmentID,
	)
	return err
}

func (repository *attachmentRepository) DeleteProjectAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, attachmentID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_attachments
			WHERE
				project_id = ?
			AND
				attachment_id = ?
        `,
		projectID,
		attachmentID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *attachmentRepository) GetProjectAttachments(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Attachment, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				A.id, A.creator_id, U.name, A.created_at, A.original_name, A.content_type, A.size
            FROM project_attachments PA
			INNER JOIN attachments A ON A.id = PA.attachment_id
			INNER JOIN users U ON U.id = A.creator_id
            WHERE PA.project_id = ?
			ORDER BY A.created_at DESC
        `,
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]attachmentDTO, 0)
	for rows.Next() {
		var dto attachmentDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorID, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toDomainArray(dtos), nil
}

func (repository *attachmentRepository) AddTaskAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, attachmentID string) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_attachments (task_id, attachment_id)
			VALUES (?, ?)
        `,
		taskID,
		attachmentID,
	)
	return err
}

func (repository *attachmentRepository) DeleteTaskAttachment(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, attachmentID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_attachments
			WHERE
				task_id = ?
			AND
				attachment_id = ?
        `,
		taskID,
		attachmentID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *attachmentRepository) GetTaskAttachments(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Attachment, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				A.id, A.creator_id, U.name, A.created_at, A.original_name, A.content_type, A.size
            FROM task_attachments TA
			INNER JOIN attachments A ON A.id = TA.attachment_id
			INNER JOIN users U ON U.id = A.creator_id
            WHERE TA.task_id = ?
			ORDER BY A.created_at DESC
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]attachmentDTO, 0)
	for rows.Next() {
		var dto attachmentDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorID, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toDomainArray(dtos), nil
}
