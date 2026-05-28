package noterepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type NoteRepository interface {
	AddProjectNote(ctx context.Context, projectId string, note noteDTO) error
	UpdateProjectNote(ctx context.Context, note noteDTO) error
	DeleteProjectNote(ctx context.Context, id string) error
	SearchProjectNotes(ctx context.Context, projectId string) ([]noteDTO, error)
}

type noteRepository struct {
	database database.Database
}

func NewNoteRepository(database database.Database) NoteRepository {
	return &noteRepository{database: database}
}

func (noteRepository *noteRepository) AddProjectNote(ctx context.Context, projectId string, note noteDTO) error {
	fmt.Println(note.CreatedAt)
	_, err := noteRepository.database.ExecContext(
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
	return err
}

func (noteRepository *noteRepository) UpdateProjectNote(ctx context.Context, note noteDTO) error {
	_, err := noteRepository.database.ExecContext(
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
	return err
}

func (noteRepository *noteRepository) DeleteProjectNote(ctx context.Context, id string) error {
	_, err := noteRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		id,
	)
	return err
}

func (noteRepository *noteRepository) SearchProjectNotes(ctx context.Context, projectId string) ([]noteDTO, error) {
	rows, err := noteRepository.database.QueryContext(
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
