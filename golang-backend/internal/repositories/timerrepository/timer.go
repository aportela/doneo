package timerrepository

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

type TimerRepository interface {
	Start(ctx context.Context, id string, userId string, summary string, startedAt int64) error
	Stop(ctx context.Context, id string, userId string, finishedAt int64) error
	Delete(ctx context.Context, id string, userId string) error
	Clear(ctx context.Context, userId string) error
	Search(ctx context.Context, userId string) ([]domain.Timer, error)
}

type timerRepository struct {
	db database.Database
}

func NewRepository(db database.Database) TimerRepository {
	return &timerRepository{db: db}
}

func (repository *timerRepository) Start(ctx context.Context, id string, userId string, summary string, startedAt int64) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO timers (id, user_id, summary, started_at, finished_at)
			VALUES (?, ?, ?, ?, NULL)
        `,
		id,
		userId,
		summary,
		startedAt,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "user_id"}
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

func (repository *timerRepository) Stop(ctx context.Context, id string, userId string, finishedAt int64) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            UPDATE timers SET
				finished_at = ?
			WHERE id = ?
			AND user_id = ?
        `,
		finishedAt,
		id,
		userId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "user_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *timerRepository) Delete(ctx context.Context, id string, userId string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM timers
			WHERE id = ?
			AND user_id = ?
        `,
		id,
		userId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "user_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *timerRepository) Clear(ctx context.Context, userId string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`
            DELETE FROM timers
			WHERE user_id = ?
        `,
		userId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE:
			// TODO
			return err
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(id)") {
				return &domain.ValidationError{Field: "id"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "user_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *timerRepository) Search(ctx context.Context, userId string) ([]domain.Timer, error) {
	rows, err := repository.db.QueryContext(ctx,
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
	dtos := make([]timerDTO, 0)
	for rows.Next() {
		var dto timerDTO
		if err := rows.Scan(
			&dto.ID, &dto.Summary, &dto.StartedAt, &dto.FinishedAt,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	return toDomainArray(dtos), nil
}
