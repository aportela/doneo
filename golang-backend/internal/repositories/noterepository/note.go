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
	UpdateProjectNote(ctx context.Context, note domain.Note) error
	DeleteProjectNote(ctx context.Context, id string) error
	GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)
	AddTaskNote(ctx context.Context, taskId string, note domain.Note) error
	UpdateTaskNote(ctx context.Context, note domain.Note) error
	DeleteTaskNote(ctx context.Context, id string) error
	GetTaskNotes(ctx context.Context, taskId string) ([]domain.Note, error)
}

type noteRepository struct {
	db database.Database
}

func NewRepository(db database.Database) NoteRepository {
	return &noteRepository{db: db}
}

func (repository *noteRepository) AddProjectNote(ctx context.Context, projectId string, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO project_notes (id, project_id, creator_id, created_at, updated_at, body)
			VALUES (?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		projectId,
		dto.CreatorId,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		// TODO: remove ?
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
			} else if strings.Contains(sqlErr.Error(), "length(creator_id)") {
				return &domain.ValidationError{Field: "creator_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *noteRepository) UpdateProjectNote(ctx context.Context, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.db.ExecContext(
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
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) DeleteProjectNote(ctx context.Context, id string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		id,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`
            SELECT
				PN.id, PN.creator_id, U.name, PN.created_at, PN.updated_at, PN.body
            FROM project_notes PN
			INNER JOIN users U ON U.id = PN.creator_id
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
			&dto.ID, &dto.CreatorId, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
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

func (repository *noteRepository) AddTaskNote(ctx context.Context, taskId string, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO task_notes (id, task_id, creator_id, created_at, updated_at, body)
			VALUES (?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		taskId,
		dto.CreatorId,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
			} else if strings.Contains(sqlErr.Error(), "length(creator_id)") {
				return &domain.ValidationError{Field: "creator_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *noteRepository) UpdateTaskNote(ctx context.Context, note domain.Note) error {
	dto := toDTO(note)
	_, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE task_notes SET
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
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) DeleteTaskNote(ctx context.Context, id string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM task_notes
			WHERE
				id = ?
        `,
		id,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *noteRepository) GetTaskNotes(ctx context.Context, taskId string) ([]domain.Note, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`
            SELECT
				TN.id, TN.creator_id, U.name, TN.created_at, TN.updated_at, TN.body
            FROM task_notes TN
			INNER JOIN users U ON U.id = TN.creator_id
            WHERE TN.task_id = ?
			ORDER BY TN.created_at DESC
        `,
		taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]noteDTO, 0)
	for rows.Next() {
		var dto noteDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorId, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
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
