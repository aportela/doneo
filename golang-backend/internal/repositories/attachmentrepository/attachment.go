package attachmentrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type AttachmentRepository interface {
	GetAttachment(ctx context.Context, attachmentId string) (domain.Attachment, error)
	AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error
	DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	GetProjectAttachments(ctx context.Context, projectId string) ([]domain.Attachment, error)
}

type attachmentRepository struct {
	database database.Database
}

func NewRepository(database database.Database) AttachmentRepository {
	return &attachmentRepository{database: database}
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

func (repository *attachmentRepository) AddProjectAttachment(ctx context.Context, projectId string, attachment domain.Attachment) error {
	dto := toDTO(attachment)
	tx, err := repository.database.Begin()
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
	_, err = tx.ExecContext(
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
		}
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		`
            INSERT INTO project_attachments (project_id, attachment_id)
			VALUES (?, ?)
        `,
		projectId,
		dto.ID,
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
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
			} else if strings.Contains(sqlErr.Error(), "length(attachmnet_id)") {
				return &domain.ValidationError{Field: "attachmnet_id"}
			}
		}
		return err
	}
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectAttachmentAdded,
		dto.UserId,
		dto.CreatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *attachmentRepository) DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error {
	tx, err := repository.database.Begin()
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
	_, err = tx.ExecContext(
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
	_, err = tx.ExecContext(
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
	userId, _ := middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectAttachmentDeleted,
		userId,
		time.Now().UnixMilli(),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
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
	attachments := make([]attachmentDTO, 0)
	for rows.Next() {
		var dto attachmentDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.OriginalName, &dto.ContentType, &dto.Size,
		); err != nil {
			return nil, err
		}
		attachments = append(attachments, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toDomainArray(attachments), nil
}
