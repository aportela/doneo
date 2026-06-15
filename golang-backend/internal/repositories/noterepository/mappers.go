package noterepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(note domain.Note) noteDTO {
	return noteDTO{
		ID:          note.ID,
		CreatorId:   note.CreatedBy.ID,
		CreatorName: note.CreatedBy.Name,
		CreatedAt:   note.CreatedAt.UnixMilli(),
		UpdatedAt:   utils.TimePtrToSQLNullInt64(note.UpdatedAt),
		Body:        note.Body,
	}
}

func toDomain(note noteDTO) domain.Note {
	return domain.Note{
		ID: note.ID,
		CreatedBy: domain.UserBase{
			ID:   note.CreatorId,
			Name: note.CreatorName,
		},
		CreatedAt: time.UnixMilli(note.CreatedAt),
		UpdatedAt: utils.SQLNullInt64ToTimePtr(note.UpdatedAt),
		Body:      note.Body,
	}
}

func toDomainArray(notes []noteDTO) []domain.Note {
	results := make([]domain.Note, 0, len(notes))
	for _, note := range notes {
		results = append(results, toDomain(note))
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
		if strings.Contains(sqlErr.Error(), "length(body)") {
			return &domain.ValidationError{Field: "body"}
		}
	}

	return err
}
