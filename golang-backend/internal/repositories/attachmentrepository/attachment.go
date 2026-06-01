package attachmentrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type AttachmentRepository interface {
	AddProjectAttachment(ctx context.Context, projectId string, attachment attachmentDTO) error
	DeleteAttachment(ctx context.Context, id string) error
}

type attachmentRepository struct {
	database database.Database
}

func NewAttachmentRepository(database database.Database) AttachmentRepository {
	return &attachmentRepository{database: database}
}

func (attachmentRepository *attachmentRepository) AddProjectAttachment(ctx context.Context, projectId string, attachment attachmentDTO) error {
	tx, err := attachmentRepository.database.Begin()
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

func (attachmentRepository *attachmentRepository) DeleteAttachment(ctx context.Context, id string) error {
	_, err := attachmentRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM attachments
			WHERE
				id = ?
        `,
		id,
	)
	return err
}
