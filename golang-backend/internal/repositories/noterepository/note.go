package noterepository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type NoteRepository interface {
	AddProjectNote(ctx context.Context, projectId string, note noteDTO) error
	UpdateProjectNote(ctx context.Context, projectId string, note noteDTO) error
	DeleteProjectNote(ctx context.Context, projectId string, id string) error
	GetProjectNotes(ctx context.Context, projectId string) ([]noteDTO, error)
}

type noteRepository struct {
	database database.Database
}

func NewRepository(database database.Database) NoteRepository {
	return &noteRepository{database: database}
}

func (repository *noteRepository) AddProjectNote(ctx context.Context, projectId string, note noteDTO) error {
	tx, err := repository.database.Begin()
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
	_, err = tx.ExecContext(
		ctx,
		`
            INSERT INTO project_notes (id, project_id, user_id, created_at, updated_at, body)
			VALUES (?, ?, ?, ?, NULL, ?)
        `,
		note.ID,
		projectId,
		note.UserId,
		note.CreatedAt,
		note.Body,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "projectId"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "userId"}
			}
		}
	}
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectNoteAdded,
		note.UserId,
		note.CreatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *noteRepository) UpdateProjectNote(ctx context.Context, projectId string, note noteDTO) error {
	tx, err := repository.database.Begin()
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
	_, err = repository.database.ExecContext(
		ctx,
		`
            UPDATE project_notes SET
				updated_at = ?,
				body = ?
			WHERE
				id = ?
        `,
		note.UpdatedAt,
		note.Body,
		note.ID,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	userId, _ := middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectNoteUpdated,
		userId,
		note.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *noteRepository) DeleteProjectNote(ctx context.Context, projectId string, id string) error {
	tx, err := repository.database.Begin()
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
	_, err = repository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		id,
	)
	userId, _ := middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectNoteDeleted,
		userId,
		time.Now().UnixMilli(),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *noteRepository) GetProjectNotes(ctx context.Context, projectId string) ([]noteDTO, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
				PN.id, PN.user_id, U.name, PN.created_at, PN.updated_at, PN.body
            FROM project_notes PN
			INNER JOIN users U ON U.id = PN.user_id
            WHERE PN.project_id = ?
			ORDER BY PN.created_at DESC
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := make([]noteDTO, 0)
	for rows.Next() {
		var note noteDTO
		if err := rows.Scan(
			&note.ID, &note.UserId, &note.UserName, &note.CreatedAt, &note.UpdatedAt, &note.Body,
		); err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notes, nil
}
