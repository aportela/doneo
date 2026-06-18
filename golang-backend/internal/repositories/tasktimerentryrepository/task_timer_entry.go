package tasktimerentryrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type TaskTimerEntryRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, taskTimerEntry domain.TaskTimerEntry) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntry domain.TaskTimerEntry) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntryID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntryID string) (domain.TaskTimerEntry, error)
	GetTaskTimerEntries(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.TaskTimerEntry, error)
}

type taskTimerEntryRepository struct{}

func NewRepository() TaskTimerEntryRepository {
	return &taskTimerEntryRepository{}
}

func (repository *taskTimerEntryRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, taskTimerEntry domain.TaskTimerEntry) error {
	dto := toDTO(taskTimerEntry)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_time_trackings
				(id, task_id, creator_id, created_at, summary, total_seconds)
			VALUES
				(?, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		taskID,
		dto.CreatorID,
		dto.CreatedAt,
		dto.Summary,
		dto.TotalSeconds,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *taskTimerEntryRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntry domain.TaskTimerEntry) error {
	dto := toDTO(taskTimerEntry)
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE task_time_trackings
			SET
				summary = ?,
				total_seconds = ?
			WHERE
				id = ?
        `,
		dto.Summary,
		dto.TotalSeconds,
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

func (repository *taskTimerEntryRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntryID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_time_trackings
			WHERE
				id = ?
        `,
		taskTimerEntryID,
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

func (repository *taskTimerEntryRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimerEntryID string) (domain.TaskTimerEntry, error) {
	var dto taskTimerEntryDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                TTT.id,
				TTT.creator_id,
				U.name AS creator_name,
				TTT.created_at,
				TTT.summary,
				TTT.total_seconds
            FROM task_time_trackings TTT
			INNER JOIN users U ON U.ID = TTT.creator_id
            WHERE
				TTT.id = ?
        `,
		taskTimerEntryID).Scan(
		&dto.ID,
		&dto.CreatorID,
		&dto.CreatorName,
		&dto.CreatedAt,
		&dto.Summary,
		&dto.TotalSeconds,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskTimerEntry{}, domain.NotFoundError
		}
		return domain.TaskTimerEntry{}, err
	}
	return toDomain(dto), err
}

func (repository *taskTimerEntryRepository) GetTaskTimerEntries(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.TaskTimerEntry, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
                TTT.id,
				TTT.creator_id,
				U.name AS creator_name,
				TTT.created_at,
				TTT.summary,
				TTT.total_seconds
            FROM task_time_trackings TTT
			INNER JOIN users U ON U.ID = TTT.creator_id
            WHERE
				TTT.task_id = ?
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]taskTimerEntryDTO, 0)
	for rows.Next() {
		var dto taskTimerEntryDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorID, &dto.CreatorName, &dto.CreatedAt, &dto.Summary, &dto.TotalSeconds,
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
