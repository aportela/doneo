package projecthistoryrespository

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

type ProjectHistoryRepository interface {
	AddProjectHistoryOperation(ctx context.Context, projectId string, operationType uint, operationDate int64, operationUserId string) error
	GetProjectHistoryOperations(ctx context.Context, projectId string) ([]projectHistoryOperationDTO, error)
}

type projectHistoryRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectHistoryRepository {
	return &projectHistoryRepository{database: database}
}

func (repository *projectHistoryRepository) AddProjectHistoryOperation(ctx context.Context, projectId string, operationType uint, operationDate int64, operationUserId string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		operationType,
		operationUserId,
		operationDate,
	)
	if err != nil {
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "project_id, operation_type, user_id, created_at"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "project_id"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "user_id"}
			}
		default:
			// TODO: return sqlErr ??? (check other repositories)
			return err
		}
	}
	return nil
}

func (repository *projectHistoryRepository) GetProjectHistoryOperations(ctx context.Context, projectId string) ([]projectHistoryOperationDTO, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
				PHO.user_id, U.name, PHO.created_at, PHO.operation_type
            FROM project_history_operations PHO
			INNER JOIN users U ON U.id = PHO.user_id
            WHERE PHO.project_id = ?
			ORDER BY PHO.created_at DESC
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	operations := make([]projectHistoryOperationDTO, 0)
	for rows.Next() {
		var operation projectHistoryOperationDTO
		if err := rows.Scan(
			&operation.UserId, &operation.UserName, &operation.CreatedAt, &operation.OperationType,
		); err != nil {
			return nil, err
		}
		operations = append(operations, operation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return operations, nil
}
