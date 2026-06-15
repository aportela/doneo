package noterepository

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type NoteRepository interface {
	AddProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, note domain.Note) error
	UpdateProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, note domain.Note) error
	DeleteProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, noteID string) error
	GetProjectNotes(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Note, error)

	AddTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, note domain.Note) error
	UpdateTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, note domain.Note) error
	DeleteTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, noteID string) error
	GetTaskNotes(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Note, error)
}

type noteRepository struct{}

func NewRepository() NoteRepository {
	return &noteRepository{}
}

func (repository *noteRepository) AddProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, note domain.Note) error {
	dto := toDTO(note)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_notes
				(id, project_id, creator_id, created_at, updated_at, body)
			VALUES
				(?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		projectID,
		dto.CreatorID,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *noteRepository) UpdateProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, note domain.Note) error {
	dto := toDTO(note)
	result, err := dbExecutor.ExecContext(
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

func (repository *noteRepository) DeleteProjectNote(ctx context.Context, dbExecutor database.DatabaseExecutor, noteID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_notes
			WHERE
				id = ?
        `,
		noteID,
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

func (repository *noteRepository) GetProjectNotes(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Note, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				PN.id, PN.creator_id, U.name, PN.created_at, PN.updated_at, PN.body
            FROM project_notes PN
			INNER JOIN users U ON U.id = PN.creator_id
            WHERE
				PN.project_id = ?
			ORDER BY
				PN.created_at DESC
        `,
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]noteDTO, 0)
	for rows.Next() {
		var dto noteDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
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

func (repository *noteRepository) AddTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, note domain.Note) error {
	dto := toDTO(note)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_notes
				(id, task_id, creator_id, created_at, updated_at, body)
			VALUES
				(?, ?, ?, ?, NULL, ?)
        `,
		dto.ID,
		taskID,
		dto.CreatorID,
		dto.CreatedAt,
		dto.Body,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *noteRepository) UpdateTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, note domain.Note) error {
	dto := toDTO(note)
	result, err := dbExecutor.ExecContext(
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

func (repository *noteRepository) DeleteTaskNote(ctx context.Context, dbExecutor database.DatabaseExecutor, noteID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_notes
			WHERE
				id = ?
        `,
		noteID,
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

func (repository *noteRepository) GetTaskNotes(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Note, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				TN.id, TN.creator_id, U.name, TN.created_at, TN.updated_at, TN.body
            FROM task_notes TN
			INNER JOIN users U ON U.id = TN.creator_id
            WHERE
				TN.task_id = ?
			ORDER BY
				TN.created_at DESC
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]noteDTO, 0)
	for rows.Next() {
		var dto noteDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt, &dto.Body,
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
