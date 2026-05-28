package noteservice

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/noterepository"
)

type NoteService interface {
	AddProjectNote(ctx context.Context, projectId string, note domain.Note) error
	UpdateProjectNote(ctx context.Context, note domain.Note) error
	DeleteProjectNote(ctx context.Context, projectId string) error
	SearchProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)
}

type noteService struct {
	repository noterepository.NoteRepository
}

func NewNoteService(repository noterepository.NoteRepository) NoteService {
	return &noteService{repository: repository}
}

func (s *noteService) AddProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	return s.repository.AddProjectNote(ctx, projectId, noterepository.DomainToDTO(note))
}

func (s *noteService) UpdateProjectNote(ctx context.Context, note domain.Note) error {
	return s.repository.UpdateProjectNote(ctx, noterepository.DomainToDTO(note))
}

func (s *noteService) DeleteProjectNote(ctx context.Context, projectId string) error {
	return s.repository.DeleteProjectNote(ctx, projectId)
}

func (s *noteService) SearchProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error) {
	notes, err := s.repository.SearchProjectNotes(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to get project permissions: %w", err)
	}
	return noterepository.DTOArrayToDomainArray(notes), nil
}
