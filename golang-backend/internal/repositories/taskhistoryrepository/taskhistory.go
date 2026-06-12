package taskhistoryrepository

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

type TaskHistoryRepository interface {
	Add(ctx context.Context, taskId string, operation domain.HistoryOperation) error
	Search(ctx context.Context, taskId string) ([]domain.HistoryOperation, error)
}

type taskHistoryRepository struct {
	database database.Database
}

func NewRepository(database database.Database) TaskHistoryRepository {
	return &taskHistoryRepository{database: database}
}

func (repository *taskHistoryRepository) Add(ctx context.Context, taskId string, operation domain.HistoryOperation) error {
	dto := toDTO(operation)
	_, err := repository.database.ExecContext(
		ctx,
		`
			INSERT INTO task_history_operations
				(id, task_id, operation_type, user_id, operation_date)
			VALUES
				(?, ?, ?, ?, ?)
		`,
		dto.ID,
		taskId,
		dto.OperationType,
		dto.UserId,
		dto.CreatedAt,
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
			return &domain.ValidationError{Field: "task_id, operation_type, user_id, operation_date"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
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

func (repository *taskHistoryRepository) Search(ctx context.Context, taskId string) ([]domain.HistoryOperation, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
				THO.id, THO.user_id, U.name, THO.operation_date, THO.operation_type
            FROM task_history_operations THO
			INNER JOIN users U ON U.id = THO.user_id
            WHERE THO.task_id = ?
			ORDER BY THO.operation_date DESC
        `,
		taskId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]taskHistoryOperationDTO, 0)
	for rows.Next() {
		var dto taskHistoryOperationDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.OperationType,
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
