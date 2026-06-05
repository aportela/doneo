package projecthistoryrepository

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
	Add(ctx context.Context, projectId string, operation domain.ProjectHistoryOperation) error
	Search(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error)
}

type projectHistoryRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectHistoryRepository {
	return &projectHistoryRepository{database: database}
}

func (repository *projectHistoryRepository) Add(ctx context.Context, projectId string, operation domain.ProjectHistoryOperation) error {
	dto := toDTO(operation)
	_, err := repository.database.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		dto.OperationType,
		dto.UserId,
		dto.CreatedAt,
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

func (repository *projectHistoryRepository) Search(ctx context.Context, projectId string) ([]domain.ProjectHistoryOperation, error) {
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
	dtos := make([]projectHistoryOperationDTO, 0)
	for rows.Next() {
		var dto projectHistoryOperationDTO
		if err := rows.Scan(
			&dto.UserId, &dto.UserName, &dto.CreatedAt, &dto.OperationType,
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
