package timerrepository

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type UserTimerRepository interface {
	StartUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string, summary string, startedAt int64) error
	StopUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string, finishedAt int64) error
	DeleteUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string) error
	ClearUserTimers(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error
	GetUserTimers(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) ([]domain.UserTimer, error)
}

type userTimerRepository struct{}

func NewRepository() UserTimerRepository {
	return &userTimerRepository{}
}

func (repository *userTimerRepository) StartUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string, summary string, startedAt int64) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO timers
				(id, user_id, summary, started_at, finished_at)
			VALUES
				(?, ?, ?, ?, NULL)
        `,
		timerID,
		userID,
		summary,
		startedAt,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *userTimerRepository) StopUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string, finishedAt int64) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            UPDATE timers
			SET
				finished_at = ?
			WHERE
				id = ?
			AND
				user_id = ?
        `,
		finishedAt,
		timerID,
		userID,
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

func (repository *userTimerRepository) DeleteUserTimer(ctx context.Context, dbExecutor database.DatabaseExecutor, timerID string, userID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM timers
			WHERE
				id = ?
			AND
				user_id = ?
        `,
		timerID,
		userID,
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

func (repository *userTimerRepository) ClearUserTimers(ctx context.Context, dbExecutor database.DatabaseExecutor, userID string) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM timers
			WHERE
				user_id = ?
        `,
		userID,
	)
	return err
}

func (repository *userTimerRepository) GetUserTimers(ctx context.Context, dbExecutor database.DatabaseExecutor, userId string) ([]domain.UserTimer, error) {
	rows, err := dbExecutor.QueryContext(ctx,
		`
			SELECT
				T.id, T.summary, T.started_at, T.finished_at
			FROM timers T
			WHERE T.user_id = ?
			ORDER BY T.started_at DESC
		`,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]userTimerDTO, 0)
	for rows.Next() {
		var dto userTimerDTO
		if err := rows.Scan(
			&dto.ID, &dto.Summary, &dto.StartedAt, &dto.FinishedAt,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return toDomainArray(dtos), nil
}
