package attachmentrepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(attachment domain.Attachment) attachmentDTO {
	return attachmentDTO{
		ID:           attachment.ID,
		CreatorID:    attachment.CreatedBy.ID,
		UserName:     attachment.CreatedBy.Name,
		CreatedAt:    attachment.CreatedAt.UnixMilli(),
		OriginalName: attachment.OriginalName,
		ContentType:  attachment.ContentType,
		Size:         attachment.Size,
	}
}

func toDomain(attachment attachmentDTO) domain.Attachment {
	return domain.Attachment{
		ID: attachment.ID,
		CreatedBy: domain.UserBase{
			ID:   attachment.CreatorID,
			Name: attachment.UserName,
		},
		CreatedAt:    time.UnixMilli(attachment.CreatedAt),
		OriginalName: attachment.OriginalName,
		ContentType:  attachment.ContentType,
		Size:         attachment.Size,
	}
}

func toDomainArray(attachments []attachmentDTO) []domain.Attachment {
	results := make([]domain.Attachment, 0, len(attachments))
	for _, attachment := range attachments {
		results = append(results, toDomain(attachment))
	}
	return results
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}

	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(original_name)") {
			return &domain.ValidationError{Field: "original_name"}
		}
	}

	return err
}
