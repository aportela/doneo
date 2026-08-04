package pagerepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type PageRepository interface {
	AddProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, page domain.Page) error
	UpdateProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, page domain.Page) error
	DeleteProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, pageID string) error
	GetProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, pageID string) (domain.Page, error)
	GetProjectPages(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Page, error)

	AddTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, page domain.Page) error
	UpdateTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, page domain.Page) error
	DeleteTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, pageID string) error
	GetTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, pageID string) (domain.Page, error)
	GetTaskPages(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Page, error)
}

type pageRepository struct{}

func NewRepository() PageRepository {
	return &pageRepository{}
}

func (repository *pageRepository) AddProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, page domain.Page) error {
	dto := toDTO(page)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_pages
				(id, project_id, title, body, creator_id, created_at, updated_at, deleted_at)
			VALUES
				(?, ?, ?, ?, ?, ?, NULL, NULL)
        `,
		dto.ID,
		projectID,
		dto.Title,
		dto.Body,
		dto.CreatorID,
		dto.CreatedAt,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *pageRepository) UpdateProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, page domain.Page) error {
	dto := toDTO(page)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE project_pages
			SET
				updated_at = ?,
				title = ?,
				body = ?
			WHERE
				id = ?
			AND
				project_id = ?
        `,
		dto.UpdatedAt,
		dto.Title,
		dto.Body,
		dto.ID,
		projectID,
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

func (repository *pageRepository) DeleteProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, pageID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_pages
			WHERE
				id = ?
			AND
				project_id = ?
        `,
		pageID,
		projectID,
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

func (repository *pageRepository) GetProjectPage(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, pageID string) (domain.Page, error) {
	var dto pageDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
			SELECT
				PP.id, PP.title, PP.body, PP.creator_id, U.name, PP.created_at, PP.updated_at
            FROM project_pages PP
			INNER JOIN users U ON U.id = PP.creator_id
            WHERE
				PP.id = ?
			AND
				PP.project_id = ?
		`,
		pageID, projectID).Scan(&dto.ID, &dto.Title, &dto.Body, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Page{}, domain.NotFoundError
		}
		return domain.Page{}, err
	}
	return toDomain(dto), err
}

func (repository *pageRepository) GetProjectPages(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.Page, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				PP.id, PP.title, PP.creator_id, U.name, PP.created_at, PP.updated_at
            FROM project_pages PP
			INNER JOIN users U ON U.id = PP.creator_id
            WHERE
				PP.project_id = ?
			ORDER BY
				PP.created_at DESC
        `,
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]pageDTO, 0)
	for rows.Next() {
		var dto pageDTO
		if err := rows.Scan(
			&dto.ID, &dto.Title, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt,
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

func (repository *pageRepository) AddTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, page domain.Page) error {
	dto := toDTO(page)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_pages
				(id, task_id, title, body, creator_id, created_at, updated_at, deleted_at)
			VALUES
				(?, ?, ?, ?, ?, ?, NULL, NULL)
        `,
		dto.ID,
		taskID,
		dto.Title,
		dto.Body,
		dto.CreatorID,
		dto.CreatedAt,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *pageRepository) UpdateTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, page domain.Page) error {
	dto := toDTO(page)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE task_pages
			SET
				updated_at = ?,
				title = ?,
				body = ?
			WHERE
				id = ?
			AND
				task_id = ?
        `,
		dto.UpdatedAt,
		dto.Title,
		dto.Body,
		dto.ID,
		taskID,
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

func (repository *pageRepository) DeleteTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, pageID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_pages
			WHERE
				id = ?
			AND
				task_id = ?
        `,
		pageID,
		taskID,
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

func (repository *pageRepository) GetTaskPage(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, pageID string) (domain.Page, error) {
	var dto pageDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
			SELECT
				TP.id, TP.title, TP.body, TP.creator_id, U.name, TP.created_at, TP.updated_at
            FROM task_pages TP
			INNER JOIN users U ON U.id = TP.creator_id
            WHERE
				TP.id = ?
			AND
				TP.task_id = ?
		`,
		pageID, taskID).Scan(&dto.ID, &dto.Title, &dto.Body, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Page{}, domain.NotFoundError
		}
		return domain.Page{}, err
	}
	return toDomain(dto), err
}

func (repository *pageRepository) GetTaskPages(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.Page, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				TP.id, TP.title, TP.creator_id, U.name, TP.created_at, TP.updated_at
            FROM task_pages TP
			INNER JOIN users U ON U.id = TP.creator_id
            WHERE
				TP.task_id = ?
			ORDER BY
				PP.created_at DESC
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]pageDTO, 0)
	for rows.Next() {
		var dto pageDTO
		if err := rows.Scan(
			&dto.ID, &dto.Title, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.UpdatedAt,
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
