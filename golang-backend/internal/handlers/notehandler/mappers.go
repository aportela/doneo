package notehandler

import (
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/utils"
)

func addRequestToDomain(request addRequest) domain.Note {
	return domain.Note{
		Body: request.Body,
		User: domain.UserBase{},
	}
}

func updateRequestToDomain(request updateRequest) domain.Note {
	return domain.Note{
		ID:   request.ID,
		Body: request.Body,
		User: domain.UserBase{
			ID:   request.User.ID,
			Name: request.User.Name,
		},
		CreatedAt: utils.MSTimestampToTime(request.CreatedAt),
	}
}

func domainToResponse(note domain.Note) NoteResponse {
	return NoteResponse{
		ID:        note.ID,
		User:      userhandler.BaseDomainToBaseResponse(note.User),
		CreatedAt: note.CreatedAt.UnixMilli(),
		UpdatedAt: utils.TimePtrToInt64Ptr(note.UpdatedAt),
		Body:      note.Body,
	}
}

func domainArrayToResponseArray(notes []domain.Note) []NoteResponse {
	notesResponse := []NoteResponse{}
	for _, note := range notes {
		notesResponse = append(notesResponse, domainToResponse(note))
	}
	return notesResponse
}

func toSearchResponse(notes []domain.Note) searchResponse {
	return searchResponse{
		ProjectPermissions: domainArrayToResponseArray(notes),
	}
}
