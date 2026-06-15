package tasktimeentryrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type TaskTimeEntryRepository interface {
	Add(ctx context.Context, taskId string, taskTimeEntry domain.TaskTimerEntry) error
	Update(ctx context.Context, taskTimeEntry domain.TaskTimerEntry) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (domain.TaskTimerEntry, error)
	GetTaskTimeEntries(ctx context.Context, taskId string) ([]domain.TaskTimerEntry, error)
}

type taskTimeEntryRepository struct {
	db database.Database
}

func NewRepository(db database.Database) TaskTimeEntryRepository {
	return &taskTimeEntryRepository{db: db}
}

func (repository *taskTimeEntryRepository) Add(ctx context.Context, taskId string, taskTimeEntry domain.TaskTimerEntry) error {
	dto := toDTO(taskTimeEntry)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO task_timer_entries
				(id, task_id, user_id, created_at, summary, total_seconds)
			VALUES
				(?, ?, ?, ?, ?, ?)
        `,
		dto.ID,
		taskId,
		dto.CreatorId,
		dto.CreatedAt,
		dto.Summary,
		dto.TotalSeconds,
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
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO: check unique
			return err
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskTimeEntryRepository) Update(ctx context.Context, taskTimeEntry domain.TaskTimerEntry) error {
	dto := toDTO(taskTimeEntry)
	_, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE task_timer_entries SET
				summary = ?,
				total_seconds = ?
			WHERE id = ?
        `,
		dto.Summary,
		dto.TotalSeconds,
		dto.ID,
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
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO: check unique
			return err
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
			} else if strings.Contains(sqlErr.Error(), "length(summary)") {
				return &domain.ValidationError{Field: "summary"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskTimeEntryRepository) Delete(ctx context.Context, id string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM task_timer_entries
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

func (repository *taskTimeEntryRepository) Get(ctx context.Context, id string) (domain.TaskTimerEntry, error) {
	var dto taskTimerEntryDTO
	err := repository.db.QueryRowContext(
		ctx,
		`
            SELECT
                TTE.id,
				TTE.user_id,
				U.name AS creator_name,
				TTE.created_at,
				TTE.summary,
				TTE.total_seconds
            FROM task_timer_entries TTE
			INNER JOIN users U ON U.ID = TTE.user_id
            WHERE TTE.id = ?
        `,
		id).Scan(
		&dto.ID,
		&dto.CreatorId,
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

func (repository *taskTimeEntryRepository) GetTaskTimeEntries(ctx context.Context, taskId string) ([]domain.TaskTimerEntry, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`
            SELECT
                TTE.id,
				TTE.user_id,
				U.name AS creator_name,
				TTE.created_at,
				TTE.summary,
				TTE.total_seconds
            FROM task_timer_entries TTE
			INNER JOIN users U ON U.ID = TTE.user_id
            WHERE TTE.task_id = ?
        `,
		taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]taskTimerEntryDTO, 0)
	for rows.Next() {
		var dto taskTimerEntryDTO
		if err := rows.Scan(
			&dto.ID, &dto.CreatorId, &dto.CreatorName, &dto.CreatedAt, &dto.Summary, &dto.TotalSeconds,
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
