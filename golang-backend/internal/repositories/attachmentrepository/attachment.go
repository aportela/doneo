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
	GetAttachment(ctx context.Context, attachmentId string) (attachmentDTO, error)
	AddProjectAttachment(ctx context.Context, projectId string, attachment attachmentDTO) error
	DeleteProjectAttachment(ctx context.Context, projectId string, attachmentId string) error
	GetProjectAttachments(ctx context.Context, projectId string) ([]attachmentDTO, error)
}

type attachmentRepository struct {
	database database.Database
}

func NewRepository(database database.Database) AttachmentRepository {
	return &attachmentRepository{database: database}
}

func (repository *attachmentRepository) GetAttachment(ctx context.Context, id string) (attachmentDTO, error) {
	var attachment attachmentDTO
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
		id).Scan(&attachment.ID, &attachment.UserId, &attachment.UserName, &attachment.CreatedAt, &attachment.OriginalName, &attachment.ContentType, &attachment.Size)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return attachmentDTO{}, domain.NotFoundError
		}
		return attachmentDTO{}, err
	}
	return attachment, err
}

func (repository *attachmentRepository) AddProjectAttachment(ctx context.Context, projectId string, attachment attachmentDTO) error {
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
		attachment.ID,
		attachment.OriginalName,
		attachment.ContentType,
		attachment.Size,
		attachment.UserId,
		attachment.CreatedAt,
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
		attachment.ID,
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
	err = tx.Commit()
	return err
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
	return err
}

func (repository *attachmentRepository) GetProjectAttachments(ctx context.Context, projectId string) ([]attachmentDTO, error) {
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
		var attachment attachmentDTO
		if err := rows.Scan(
			&attachment.ID, &attachment.UserId, &attachment.UserName, &attachment.CreatedAt, &attachment.OriginalName, &attachment.ContentType, &attachment.Size,
		); err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return attachments, nil
}
