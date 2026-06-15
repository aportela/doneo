package tagrepository

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type TagRepository interface {
	AddTaskTag(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, tag string) error
	DeleteTaskTags(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error
	GetTaskTags(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]string, error)
}

type tagRepository struct {
}

func NewRepository() TagRepository {
	return &tagRepository{}
}

func (repository *tagRepository) AddTaskTag(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string, tag string) error {
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_tags
				(task_id, tag)
			VALUES
				(?, ?)
        `,
		taskID,
		tag,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *tagRepository) DeleteTaskTags(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM task_tags
			WHERE
				task_id = ?
        `,
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

func (repository *tagRepository) GetTaskTags(ctx context.Context, dbExecutor database.DatabaseExecutor, taskID string) ([]string, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				TT.tag
            FROM task_tags TT
            WHERE TT.task_id = ?
			ORDER BY TT.tag
        `,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tags := make([]string, 0)
	for rows.Next() {
		var tag string
		if err := rows.Scan(
			&tag,
		); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}
