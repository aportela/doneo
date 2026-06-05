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
	AddProjectNote(ctx context.Context, projectId string, note domain.Note) error
	UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) error
	DeleteProjectNote(ctx context.Context, projectId string, id string) error
	GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)
}

type noteRepository struct {
	database database.Database
}

func NewRepository(database database.Database) NoteRepository {
	return &noteRepository{database: database}
}

func (repository *noteRepository) AddProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_notes (id, project_id, user_id, created_at, updated_at, body)
			VALUES (?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		projectId,
		dto.UserId,
		dto.CreatedAt,
		dto.Body,
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
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *noteRepository) UpdateProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.database.ExecContext(
		ctx,
		`
            UPDATE project_notes SET
				updated_at = ?,
				body = ?
			WHERE
				id = ?
        `,
		dto.UpdatedAt,
		dto.Body,
		dto.ID,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) DeleteProjectNote(ctx context.Context, projectId string, id string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		id,
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error) {
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
	dtos := make([]noteDTO, 0)
	for rows.Next() {
		var dto noteDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toDomainArray(dtos), nil
}
