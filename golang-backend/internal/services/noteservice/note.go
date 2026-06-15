package noteservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type NoteService interface {
	AddProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error)
	UpdateProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error)
	DeleteProjectNote(ctx context.Context, projectID string, noteID string) error
	GetProjectNote(ctx context.Context, projectID string, noteID string) (domain.Note, error)
	GetProjectNotes(ctx context.Context, projectID string) ([]domain.Note, error)

	AddTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error)
	UpdateTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error)
	DeleteTaskNote(ctx context.Context, projectID string, taskID string, noteID string) error
	GetTaskNote(ctx context.Context, projectID string, noteID string) (domain.Note, error)
	GetTaskNotes(ctx context.Context, projectID string, taskID string) ([]domain.Note, error)
}

type noteService struct {
	database                database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	noteRepository          noterepository.NoteRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, repository noterepository.NoteRepository) NoteService {
	return &noteService{database: db, historyOperationService: historyOperationService, authorizationService: authorizationService, noteRepository: repository}
}

func (service *noteService) withProjectUpdatePermission(ctx context.Context, projectID string, action func(userID string) error) error {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	if err := service.authorizationService.RequireProjectUpdatePermission(ctx, currentContextUserID, projectID); err != nil {
		return err
	}

	return action(currentContextUserID)
}

func (service *noteService) AddProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		note.ID = utils.UUID()
		note.CreatedBy.ID = currentUserID
		note.CreatedAt = time.Now()

		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {
			if err := service.noteRepository.AddProjectNote(ctx, tx, projectID, note); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     note.CreatedAt,
					OperationType: domain.EventProjectNoteAdded,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return note, err
}

func (service *noteService) UpdateProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		note.UpdatedAt = utils.CurrentTimePtr()

		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {

			if err := service.noteRepository.UpdateProjectNote(ctx, tx, note); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: currentUserID,
					},
					CreatedAt:     *note.UpdatedAt,
					OperationType: domain.EventProjectNoteUpdated,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return note, err
}

func (service *noteService) DeleteProjectNote(ctx context.Context, projectID string, noteID string) error {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {
			if err := service.noteRepository.DeleteProjectNote(ctx, tx, noteID); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: currentUserID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectNoteDeleted,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return err
}

func (service *noteService) GetProjectNote(ctx context.Context, projectID string, noteID string) (domain.Note, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return domain.Note{}, err
	}
	note, err := service.noteRepository.GetProjectNote(ctx, service.database, noteID)
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) GetProjectNotes(ctx context.Context, projectID string) ([]domain.Note, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	notes, err := service.noteRepository.GetProjectNotes(ctx, service.database, projectID)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get project notes: %w", err)
	}
	return notes, nil
}

func (service *noteService) AddTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		note.ID = utils.UUID()
		note.CreatedBy.ID = currentUserID
		note.CreatedAt = time.Now()

		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {
			if err := service.noteRepository.AddTaskNote(ctx, tx, taskID, note); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: currentUserID},
					CreatedAt:     note.CreatedAt,
					OperationType: domain.EventTaskNoteAdded,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return note, err
}

func (service *noteService) UpdateTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		note.UpdatedAt = utils.CurrentTimePtr()

		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {

			if err := service.noteRepository.UpdateTaskNote(ctx, tx, note); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: currentUserID,
					},
					CreatedAt:     *note.UpdatedAt,
					OperationType: domain.EventTaskNoteUpdated,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return note, err
}

func (service *noteService) DeleteTaskNote(ctx context.Context, projectID string, taskID string, noteID string) error {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		return database.WithTx(ctx, service.database, func(tx *sql.Tx) error {
			if err := service.noteRepository.DeleteTaskNote(ctx, tx, noteID); err != nil {
				return err
			}

			if _, err := service.historyOperationService.AddTaskHistoryOperation(
				ctx,
				tx,
				projectID,
				taskID,
				domain.HistoryOperation{
					ID: utils.UUID(),
					CreatedBy: domain.UserBase{
						ID: currentUserID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskNoteDeleted,
				},
			); err != nil {
				return err
			}

			return nil
		})
	})

	return err
}

func (service *noteService) GetTaskNote(ctx context.Context, projectID string, noteID string) (domain.Note, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return domain.Note{}, err
	}
	note, err := service.noteRepository.GetTaskNote(ctx, service.database, noteID)
	if err != nil {
		return domain.Note{}, fmt.Errorf("[NoteService] failed to get task note: %w", err)
	}
	return note, nil
}

func (service *noteService) GetTaskNotes(ctx context.Context, projectID string, taskID string) ([]domain.Note, error) {
	currentContextUserID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	if err := service.authorizationService.RequireProjectViewPermission(ctx, currentContextUserID, projectID); err != nil {
		return nil, err
	}
	notes, err := service.noteRepository.GetTaskNotes(ctx, service.database, taskID)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get task notes: %w", err)
	}
	return notes, nil
}
