package attachmentrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type AttachmentRepository interface {
	AddAttachment(ctx context.Context, attachment domain.Attachment) error
	DeleteAttachment(ctx context.Context, attachmentId string) error
	GetAttachment(ctx context.Context, attachmentId string) (domain.Attachment, error)
	AddProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error)
}

type attachmentRepository struct {
	database database.Database
}

func NewRepository(database database.Database) AttachmentRepository {
	return &attachmentRepository{database: database}
}

func (repository *attachmentRepository) AddAttachment(ctx context.Context, attachment domain.Attachment) error {
	dto := toDTO(attachment)
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO attachments (id, original_name, content_type, size, user_id, created_at)
			VALUES (?, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		dto.OriginalName,
		dto.ContentType,
		dto.Size,
		dto.UserId,
		dto.CreatedAt,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "userId"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *attachmentRepository) DeleteAttachment(ctx context.Context, attachmentId string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM attachments
			WHERE
				id = ?
        `,
		attachmentId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *attachmentRepository) GetAttachment(ctx context.Context, id string) (domain.Attachment, error) {
	var dto attachmentDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
				A.id, A.user_id, U.name, A.created_at, A.original_name, A.content_type, A.size
            FROM project_attachments PA
			INNER JOIN attachments A ON A.id = PA.attachment_id
			INNER JOIN users U ON U.id = A.user_id
            WHERE A.id = ?
        `,
		id).Scan(&dto.ID, &dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Attachment{}, domain.NotFoundError
		}
		return domain.Attachment{}, err
	}
	return toDomain(dto), err
}

func (repository *attachmentRepository) AddProjectAttachment(ctx context.Context, projectId string, attachmentId string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_attachments (project_id, attachment_id)
			VALUES (?, ?)
        `,
		projectId,
		attachmentId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "project_id, attachment_id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
			} else if strings.Contains(sqlErr.Error(), "length(attachment_id)") {
				return &domain.ValidationError{Field: "attachment_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *attachmentRepository) DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_attachments
			WHERE
				project_id = ?
			AND
				attachment_id = ?
        `,
		projectId,
		attachmentId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *attachmentRepository) GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
				A.id, A.user_id, U.name, A.created_at, A.original_name, A.content_type, A.size
            FROM project_attachments PA
			INNER JOIN attachments A ON A.id = PA.attachment_id
			INNER JOIN users U ON U.id = A.user_id
            WHERE PA.project_id = ?
			ORDER BY A.created_at DESC
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]attachmentDTO, 0)
	for rows.Next() {
		var dto attachmentDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size,
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
