package noteservice

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/noterepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/historyoperationservice"
	"github.com/aportela/doneo/internal/utils"
)

type NoteService interface {
	AddProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error)
	UpdateProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error)
	DeleteProjectNote(ctx context.Context, projectID string, noteID string) error
	GetProjectNotes(ctx context.Context, projectID string) ([]domain.Note, error)

	AddTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error)
	UpdateTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error)
	DeleteTaskNote(ctx context.Context, projectID string, taskID string, noteID string) error
	GetTaskNotes(ctx context.Context, projectID string, taskID string) ([]domain.Note, error)
}

type noteService struct {
	db                      database.Database
	authorizationService    authorizationservice.AuthorizationService
	historyOperationService historyoperationservice.HistoryOperationService
	noteRepository          noterepository.NoteRepository
}

func NewService(db database.Database, authorizationService authorizationservice.AuthorizationService, historyOperationService historyoperationservice.HistoryOperationService, repository noterepository.NoteRepository) NoteService {
	return &noteService{db: db, historyOperationService: historyOperationService, authorizationService: authorizationService, noteRepository: repository}
}

func (service *noteService) AddProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error) {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return domain.Note{}, err
	} else {
		note.ID = utils.UUID()
		note.CreatedBy.ID = contextUser.ID
		note.CreatedAt = time.Now()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
			if err := service.noteRepository.AddProjectNote(ctx, tx, projectID, note); err != nil {
				return err
			}
			if _, err := service.historyOperationService.AddProjectHistoryOperation(
				ctx,
				tx,
				projectID,
				domain.HistoryOperation{
					ID:            utils.UUID(),
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     note.CreatedAt,
					OperationType: domain.EventProjectNoteAdded,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Note{}, err
		}
		return note, nil
	}
}

func (service *noteService) UpdateProjectNote(ctx context.Context, projectID string, note domain.Note) (domain.Note, error) {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return domain.Note{}, err
	} else {
		note.UpdatedAt = utils.CurrentTimePtr()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
						ID: contextUser.ID,
					},
					CreatedAt:     *note.UpdatedAt,
					OperationType: domain.EventProjectNoteUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Note{}, err
		}
		return note, nil
	}
}

func (service *noteService) DeleteProjectNote(ctx context.Context, projectID string, noteID string) error {
	if contextUser, err := service.authorizationService.RequireProjectUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
						ID: contextUser.ID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventProjectNoteDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *noteService) GetProjectNotes(ctx context.Context, projectID string) ([]domain.Note, error) {
	if _, err := service.authorizationService.RequireProjectViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if notes, err := service.noteRepository.GetProjectNotes(ctx, service.db, projectID); err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get project notes: %w", err)
	} else {
		return notes, nil
	}
}

func (service *noteService) AddTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return domain.Note{}, err
	} else {
		note.ID = utils.UUID()
		note.CreatedBy.ID = contextUser.ID
		note.CreatedAt = time.Now()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
					CreatedBy:     domain.UserBase{ID: contextUser.ID},
					CreatedAt:     note.CreatedAt,
					OperationType: domain.EventTaskNoteAdded,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Note{}, err
		}
		return note, nil
	}
}

func (service *noteService) UpdateTaskNote(ctx context.Context, projectID string, taskID string, note domain.Note) (domain.Note, error) {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return domain.Note{}, err
	} else {
		note.UpdatedAt = utils.CurrentTimePtr()
		if err := database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
						ID: contextUser.ID,
					},
					CreatedAt:     *note.UpdatedAt,
					OperationType: domain.EventTaskNoteUpdated,
				},
			); err != nil {
				return err
			}
			return nil
		}); err != nil {
			return domain.Note{}, err
		}
		return note, nil
	}
}

func (service *noteService) DeleteTaskNote(ctx context.Context, projectID string, taskID string, noteID string) error {
	if contextUser, err := service.authorizationService.RequireTaskUpdatePermission(ctx, projectID); err != nil {
		return err
	} else {
		return database.WithTx(ctx, service.db, func(tx *sql.Tx) error {
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
						ID: contextUser.ID,
					},
					CreatedAt:     time.Now(),
					OperationType: domain.EventTaskNoteDeleted,
				},
			); err != nil {
				return err
			}
			return nil
		})
	}
}

func (service *noteService) GetTaskNotes(ctx context.Context, projectID string, taskID string) ([]domain.Note, error) {
	if _, err := service.authorizationService.RequireTaskViewPermission(ctx, projectID); err != nil {
		return nil, err
	}
	if notes, err := service.noteRepository.GetTaskNotes(ctx, service.db, taskID); err != nil {
		return nil, fmt.Errorf("[NoteService] failed to get task notes: %w", err)
	} else {
		return notes, nil
	}
}
