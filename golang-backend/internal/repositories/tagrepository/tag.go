package tagrepository

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

type TagRepository interface {
	AddTaskTag(ctx context.Context, taskId string, tag string) error
	DeleteTaskTags(ctx context.Context, taskId string) error
	GetTaskTags(ctx context.Context, taskId string) ([]string, error)
}

type tagRepository struct {
	db database.Database
}

func NewRepository(db database.Database) TagRepository {
	return &tagRepository{db: db}
}

func (repository *tagRepository) AddTaskTag(ctx context.Context, taskId string, tag string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO task_tags (task_id, tag)
			VALUES (?, ?)
        `,
		taskId,
		tag,
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
			return &domain.ValidationError{Field: "task_id, tag"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *tagRepository) DeleteTaskTags(ctx context.Context, taskId string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM task_tags
			WHERE
				task_id = ?
        `,
		taskId,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *tagRepository) GetTaskTags(ctx context.Context, taskId string) ([]string, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`
            SELECT
				TT.tag
            FROM task_tags TT
            WHERE TT.task_id = ?
			ORDER BY TT.tag
        `,
		taskId)
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
