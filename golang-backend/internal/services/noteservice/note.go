package noteservice

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type NoteService interface {
	AddProjectNote(ctx context.Context, projectId string, note domain.Note) (domain.Note, error)
	UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) (domain.Note, error)
	DeleteProjectNote(ctx context.Context, projectId string, noteId string) error
	GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)
	AddTaskNote(ctx context.Context, projectId string, taskId string, note domain.Note) (domain.Note, error)
	UpdateTaskNote(ctx context.Context, projectId string, taskId string, note domain.Note) (domain.Note, error)
	DeleteTaskNote(ctx context.Context, projectId string, taskId string, noteId string) error
	GetTaskNotes(ctx context.Context, taskId string) ([]domain.Note, error)
}

type noteService struct {
	database                database.Database
	historyOperationService historyoperationservice.HistoryOperationService
	repository              noterepository.NoteRepository
}

func NewService(database database.Database, historyOperationService historyoperationservice.HistoryOperationService, repository noterepository.NoteRepository) NoteService {
	return &noteService{database: database, historyOperationService: historyOperationService, repository: repository}
}

func (service *noteService) AddProjectNote(ctx context.Context, projectId string, note domain.Note) (domain.Note, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Note{}, err
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
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	note.ID = utils.UUID()
	note.User.ID = currentUserId
	note.CreatedAt = time.Now()
	err = service.repository.AddProjectNote(ctx, projectId, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: note.CreatedAt, OperationType: domain.EventProjectNoteAdded})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) (domain.Note, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Note{}, err
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
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	note.UpdatedAt = utils.CurrentTimePtr()
	err = service.repository.UpdateProjectNote(ctx, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectNoteUpdated})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
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
		return fmt.Errorf("[NoteService] user ID not found in context")
	}
	err = service.repository.DeleteProjectNote(ctx, noteId)
	if err != nil {
		return err
	}
	_, err = service.historyOperationService.AddProjectHistoryOperation(ctx, projectId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventProjectNoteDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *noteService) GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error) {
	notes, err := service.repository.GetProjectNotes(ctx, projectId)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get project notes: %w", err)
	}
	return notes, nil
}

func (service *noteService) AddTaskNote(ctx context.Context, projectId string, taskId string, note domain.Note) (domain.Note, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Note{}, err
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
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	note.ID = utils.UUID()
	note.User.ID = currentUserId
	note.CreatedAt = time.Now()
	err = service.repository.AddTaskNote(ctx, taskId, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: note.CreatedAt, OperationType: domain.EventTaskNoteAdded})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) UpdateTaskNote(ctx context.Context, projectId string, taskId string, note domain.Note) (domain.Note, error) {
	tx, err := service.database.Begin()
	if err != nil {
		return domain.Note{}, err
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
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	note.UpdatedAt = utils.CurrentTimePtr()
	err = service.repository.UpdateTaskNote(ctx, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskNoteUpdated})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) DeleteTaskNote(ctx context.Context, projectId string, taskId string, noteId string) error {
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
		return fmt.Errorf("[NoteService] user ID not found in context")
	}
	err = service.repository.DeleteTaskNote(ctx, noteId)
	if err != nil {
		return err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, projectId, taskId, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskNoteDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *noteService) GetTaskNotes(ctx context.Context, taskId string) ([]domain.Note, error) {
	notes, err := service.repository.GetTaskNotes(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get task notes: %w", err)
	}
	return notes, nil
}
