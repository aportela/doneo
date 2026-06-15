package noterepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type NoteRepository interface {
	AddProjectNote(ctx context.Context, projectId string, note domain.Note) error
	UpdateProjectNote(ctx context.Context, note domain.Note) error
	DeleteProjectNote(ctx context.Context, noteId string) error
	GetProjectNote(ctx context.Context, noteId string) (domain.Note, error)
	GetProjectNotes(ctx context.Context, projectId string) ([]domain.Note, error)

	AddTaskNote(ctx context.Context, taskId string, note domain.Note) error
	UpdateTaskNote(ctx context.Context, note domain.Note) error
	DeleteTaskNote(ctx context.Context, noteId string) error
	GetTaskNote(ctx context.Context, noteId string) (domain.Note, error)
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
            INSERT INTO project_notes
				(id, project_id, creator_id, created_at, updated_at, body)
			VALUES
				(?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		projectId,
		dto.CreatorId,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *noteRepository) UpdateProjectNote(ctx context.Context, note domain.Note) error {
	dto := toDTO(note)
	result, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE project_notes
			SET
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
		return mapSQLiteError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *noteRepository) DeleteProjectNote(ctx context.Context, id string) error {
	result, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		id,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *noteRepository) GetProjectNote(ctx context.Context, noteId string) (domain.Note, error) {
	var dto noteDTO
	err := repository.db.QueryRowContext(
		ctx,
		`
            SELECT
				PN.id, PN.creator_id, U.name, PN.created_at, PN.updated_at, PN.body
            FROM project_notes PN
			INNER JOIN users U ON U.id = PN.creator_id
            WHERE PN.id = ?
        `,
		noteId).Scan(
		&dto.ID, &dto.CreatorId, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Note{}, domain.NotFoundError
		}
		return domain.Note{}, err
	}
	return toDomain(dto), err
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
            INSERT INTO task_notes
				(id, task_id, creator_id, created_at, updated_at, body)
			VALUES
				(?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		taskId,
		dto.CreatorId,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *noteRepository) UpdateTaskNote(ctx context.Context, note domain.Note) error {
	dto := toDTO(note)
	result, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE task_notes
			SET
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
		return mapSQLiteError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *noteRepository) DeleteTaskNote(ctx context.Context, id string) error {
	result, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM task_notes
			WHERE
				id = ?
        `,
		id,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *noteRepository) GetTaskNote(ctx context.Context, noteId string) (domain.Note, error) {
	var dto noteDTO
	err := repository.db.QueryRowContext(
		ctx,
		`
            SELECT
				TN.id, TN.creator_id, U.name, TN.created_at, TN.updated_at, TN.body
            FROM task_notes TN
			INNER JOIN users U ON U.id = TN.creator_id
            WHERE TN.id = ?
        `,
		noteId).Scan(
		&dto.ID, &dto.CreatorId, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Note{}, domain.NotFoundError
		}
		return domain.Note{}, err
	}
	return toDomain(dto), err
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
