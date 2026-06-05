package noteservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrepository"
	"github.com/aportela/doneo/internal/utils"
)

type NoteService interface {
	AddProjectNote(ctx context.Context, projectId string, note domain.Note) error
	UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) error
	DeleteProjectNote(ctx context.Context, projectId string, noteId string) error
	GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)
}

type noteService struct {
	database   database.Database
	repository noterepository.NoteRepository
}

func NewService(database database.Database, repository noterepository.NoteRepository) NoteService {
	return &noteService{database: database, repository: repository}
}

func (service *noteService) AddProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	tx, err := service.database.Begin()
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
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}
	err = service.repository.AddProjectNote(ctx, projectId, note)
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: note.CreatedAt, OperationType: domain.EventProjectNoteAdded})
	return tx.Commit()
}

func (service *noteService) UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	tx, err := service.database.Begin()
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
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}
	err = service.repository.UpdateProjectNote(ctx, projectId, note)
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectNoteUpdated})
	return tx.Commit()
}

func (service *noteService) DeleteProjectNote(ctx context.Context, projectId string, noteId string) error {
	tx, err := service.database.Begin()
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
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user ID not found in context")
	}
	err = service.repository.DeleteProjectNote(ctx, projectId, noteId)
	if err != nil {
		return err
	}
	err = projecthistoryrepository.NewRepository(service.database).Add(ctx, projectId, domain.ProjectHistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectNoteDeleted})
	return tx.Commit()
}

func (service *noteService) GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error) {
	notes, err := service.repository.GetProjectNotes(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[ProjectTypeService] failed to get project permissions: %w", err)
	}
	return notes, nil
}
