package tasktimetrackingrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type TaskTimeTrackingRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, taskTimeTracking domain.TaskTimeTracking) error
	Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTracking domain.TaskTimeTracking) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTrackingID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTrackingID string) (domain.TaskTimeTracking, error)
	GetTaskTimeTrackings(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.TaskTimeTracking, error)
}

type taskTimeTrackingRepository struct{}

func NewRepository() TaskTimeTrackingRepository {
	return &taskTimeTrackingRepository{}
}

func (repository *taskTimeTrackingRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, taskTimeTracking domain.TaskTimeTracking) error {
	dto := toDTO(taskTimeTracking)
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

func (repository *taskTimeTrackingRepository) Update(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTracking domain.TaskTimeTracking) error {
	dto := toDTO(taskTimeTracking)
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

func (repository *taskTimeTrackingRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTrackingID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_time_trackings
			WHERE
				id = ?
        `,
		taskTimeTrackingID,
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

func (repository *taskTimeTrackingRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, taskTimeTrackingID string) (domain.TaskTimeTracking, error) {
	var dto taskTimeTrackingDTO
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
		taskTimeTrackingID).Scan(
		&dto.ID,
		&dto.CreatorID,
		&dto.CreatorName,
		&dto.CreatedAt,
		&dto.Summary,
		&dto.TotalSeconds,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.TaskTimeTracking{}, domain.NotFoundError
		}
		return domain.TaskTimeTracking{}, err
	}
	return toDomain(dto), err
}

func (repository *taskTimeTrackingRepository) GetTaskTimeTrackings(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]domain.TaskTimeTracking, error) {
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
			ORDER BY
				TTT.created_at DESC
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]taskTimeTrackingDTO, 0)
	for rows.Next() {
		var dto taskTimeTrackingDTO
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
