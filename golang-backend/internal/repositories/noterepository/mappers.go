package noterepository

import (
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
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
