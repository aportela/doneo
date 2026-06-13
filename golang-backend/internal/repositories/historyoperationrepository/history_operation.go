package historyoperationrepository

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

type HistoryOperationRepository interface {
	AddProjectHistoryOperation(ctx context.Context, projectId string, operation domain.HistoryOperation) error
	SearchProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.HistoryOperation, error)
	AddTaskOperation(ctx context.Context, projectId string, taskId string, operation domain.HistoryOperation) error
	SearchTaskHistoryOperations(ctx context.Context, taskId string) ([]domain.HistoryOperation, error)
}

type historyOperationRepository struct {
	database database.Database
}

func NewRepository(database database.Database) HistoryOperationRepository {
	return &historyOperationRepository{database: database}
}

func (repository *historyOperationRepository) AddProjectHistoryOperation(ctx context.Context, projectId string, operation domain.HistoryOperation) error {
	dto := toDTO(operation)
	_, err := repository.database.ExecContext(
		ctx,
		`
			INSERT INTO history_operations
				(id, project_id, task_id, operation_type, user_id, operation_date)
			VALUES
				(?, ?, NULL, ?, ?, ?)
		`,
		dto.ID,
		projectId,
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
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
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

func (repository *historyOperationRepository) SearchProjectHistoryOperations(ctx context.Context, projectId string) ([]domain.HistoryOperation, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
				PHO.id, PHO.user_id, U.name, PHO.operation_date, PHO.operation_type
            FROM history_operations PHO
			INNER JOIN users U ON U.id = PHO.user_id
            WHERE PHO.project_id = ?
			ORDER BY PHO.operation_date DESC
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]historyOperationDTO, 0)
	for rows.Next() {
		var dto historyOperationDTO
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

func (repository *historyOperationRepository) AddTaskOperation(ctx context.Context, projectId string, taskId string, operation domain.HistoryOperation) error {
	dto := toDTO(operation)
	_, err := repository.database.ExecContext(
		ctx,
		`
			INSERT INTO history_operations
				(id, project_id, task_id, operation_type, user_id, operation_date)
			VALUES
				(?, ?, ?, ?, ?, ?)
		`,
		dto.ID,
		projectId,
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
			return &domain.ValidationError{Field: "id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
			} else if strings.Contains(sqlErr.Error(), "length(task_id)") {
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

func (repository *historyOperationRepository) SearchTaskHistoryOperations(ctx context.Context, taskId string) ([]domain.HistoryOperation, error) {
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
	dtos := make([]historyOperationDTO, 0)
	for rows.Next() {
		var dto historyOperationDTO
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
