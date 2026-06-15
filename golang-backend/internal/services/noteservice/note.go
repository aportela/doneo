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
	historyOperationService historyoperationservice.HistoryOperationService
	authorizationService    authorizationservice.AuthorizationService
	repository              noterepository.NoteRepository
}

func NewService(db database.Database, historyOperationService historyoperationservice.HistoryOperationService, authorizationService authorizationservice.AuthorizationService, repository noterepository.NoteRepository) NoteService {
	return &noteService{database: db, historyOperationService: historyOperationService, authorizationService: authorizationService, repository: repository}
}

func (s *noteService) withTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := s.database.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	err = fn(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (s *noteService) withProjectUpdatePermission(
	ctx context.Context,
	projectID string,
	action func(userID string) error,
) error {
	userID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	if err := s.authorizationService.RequireProjectUpdatePermission(ctx, userID, projectID); err != nil {
		return err
	}

	return action(userID)
}

func (s *noteService) withProjectViewPermission(
	ctx context.Context,
	projectID string,
	action func(userID string) error,
) error {
	userID, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("user not found in context")
	}

	if err := s.authorizationService.RequireProjectViewPermission(ctx, userID, projectID); err != nil {
		return err
	}

	return action(userID)
}

func (service *noteService) AddProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error) {
	err := service.withProjectUpdatePermission(ctx, projectID, func(currentUserID string) error {
		// TODO: mutation could be cause problems (use copy) ?
		note.ID = utils.UUID()
		note.CreatedBy.ID = currentUserID
		note.CreatedAt = time.Now()

		return service.withTx(ctx, func(tx *sql.Tx) error {
			if err := service.repository.AddProjectNote(ctx, tx, projectID, note); err != nil {
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
		// TODO: mutation could be cause problems (use copy) ?
		note.UpdatedAt = utils.CurrentTimePtr()

		return service.withTx(ctx, func(tx *sql.Tx) error {

			if err := service.repository.UpdateProjectNote(ctx, tx, note); err != nil {
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
		return service.withTx(ctx, func(tx *sql.Tx) error {
			if err := service.repository.DeleteProjectNote(ctx, tx, noteID); err != nil {
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
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireProjectViewPermission(ctx, currentUserId, projectID)
	if err != nil {
		return domain.Note{}, err
	}
	note, err := service.repository.GetProjectNote(ctx, service.database, noteID)
	if err != nil {
		return domain.Note{}, fmt.Errorf("[NoteService] failed to get project note: %w", err)
	}
	return note, nil
}

func (service *noteService) GetProjectNotes(ctx context.Context, projectID string) ([]domain.Note, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireProjectViewPermission(ctx, currentUserId, projectID)
	if err != nil {
		return nil, err
	}
	notes, err := service.repository.GetProjectNotes(ctx, service.database, projectID)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get project notes: %w", err)
	}
	return notes, nil
}

func (service *noteService) AddTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireTaskUpdatePermission(ctx, currentUserId, projectID)
	if err != nil {
		return domain.Note{}, err
	}
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
	note.ID = utils.UUID()
	note.CreatedBy.ID = currentUserId
	note.CreatedAt = time.Now()
	err = service.repository.AddTaskNote(ctx, tx, taskID, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, tx, projectID, taskID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: note.CreatedAt, OperationType: domain.EventTaskNoteAdded})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) UpdateTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireTaskUpdatePermission(ctx, currentUserId, projectID)
	if err != nil {
		return domain.Note{}, err
	}
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
	note.UpdatedAt = utils.CurrentTimePtr()
	err = service.repository.UpdateTaskNote(ctx, tx, note)
	if err != nil {
		return domain.Note{}, err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, tx, projectID, taskID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: *note.UpdatedAt, OperationType: domain.EventTaskNoteUpdated})
	if err != nil {
		return domain.Note{}, err
	}
	err = tx.Commit()
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (service *noteService) DeleteTaskNote(ctx context.Context, projectID string, taskID string, noteID string) error {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireTaskUpdatePermission(ctx, currentUserId, projectID)
	if err != nil {
		return err
	}
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
	err = service.repository.DeleteTaskNote(ctx, tx, noteID)
	if err != nil {
		return err
	}
	_, err = service.historyOperationService.AddTaskHistoryOperation(ctx, tx, projectID, taskID, domain.HistoryOperation{ID: utils.UUID(), CreatedBy: domain.UserBase{ID: currentUserId}, CreatedAt: time.Now(), OperationType: domain.EventTaskNoteDeleted})
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (service *noteService) GetTaskNote(ctx context.Context, projectID string, noteID string) (domain.Note, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return domain.Note{}, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireTaskViewPermission(ctx, currentUserId, projectID)
	if err != nil {
		return domain.Note{}, err
	}
	note, err := service.repository.GetTaskNote(ctx, service.database, noteID)
	if err != nil {
		return domain.Note{}, fmt.Errorf("[NoteService] failed to get task note: %w", err)
	}
	return note, nil
}

func (service *noteService) GetTaskNotes(ctx context.Context, projectID string, noteID string) ([]domain.Note, error) {
	currentUserId, ok := middlewares.GetUserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("[NoteService] user ID not found in context")
	}
	err := service.authorizationService.RequireTaskViewPermission(ctx, currentUserId, projectID)
	if err != nil {
		return nil, err
	}
	notes, err := service.repository.GetTaskNotes(ctx, service.database, noteID)
	if err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get task notes: %w", err)
	}
	return notes, nil
}
