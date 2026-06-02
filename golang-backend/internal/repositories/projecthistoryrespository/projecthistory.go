package projecthistoryrespository

import (
	"context"

	"github.com/aportela/doneo/internal/database"
)

type ProjectHistoryRepository interface {
	GetProjectHistoryOperations(ctx context.Context, projectId string) ([]projectHistoryOperationDTO, error)
}

type projectHistoryRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectHistoryRepository {
	return &projectHistoryRepository{database: database}
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
