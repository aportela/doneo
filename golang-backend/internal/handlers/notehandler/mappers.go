package notehandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToDomain(request addRequest) domain.Note {
	return domain.Note{
		Body: request.Body,
	}
}

func updateRequestToDomain(request updateRequest) domain.Note {
	return domain.Note{
		ID:   request.ID,
		Body: request.Body,
	}
}

func domainToResponse(note domain.Note) noteResponse {
	return noteResponse{
		ID:        note.ID,
		CreatedBy: userhandler.BaseDomainToBaseResponse(note.CreatedBy),
		CreatedAt: note.CreatedAt.UnixMilli(),
		UpdatedAt: utils.TimePtrToInt64Ptr(note.UpdatedAt),
		Body:      note.Body,
	}
}

func domainArrayToResponseArray(notes []domain.Note) []noteResponse {
	notesResponse := []noteResponse{}
	for _, note := range notes {
		notesResponse = append(notesResponse, domainToResponse(note))
	}
	return notesResponse
}

func toSearchResponse(notes []domain.Note) searchResponse {
	return searchResponse{
		Notes: domainArrayToResponseArray(notes),
	}
}
