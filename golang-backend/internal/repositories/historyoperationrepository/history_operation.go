package historyoperationrepository

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
)

type HistoryOperationRepository interface {
	addHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, taskID *string, operation domain.HistoryOperation) error
	searchHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, scope scope, scopeID string) ([]domain.HistoryOperation, error)

	AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, operation domain.HistoryOperation) error
	SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectID string) ([]domain.HistoryOperation, error)
	AddTaskHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, taskID string, operation domain.HistoryOperation) error
	SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskID string) ([]domain.HistoryOperation, error)
}

type historyOperationRepository struct{}

func NewRepository() HistoryOperationRepository {
	return &historyOperationRepository{}
}

func (repository *historyOperationRepository) addHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, taskID *string, operation domain.HistoryOperation) error {
	dto := toDTO(operation)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
			INSERT INTO history_operations
				(id, project_id, task_id, operation_type, operation_user_id, operation_date)
			VALUES
				(?, ?, ?, ?, ?, ?)
		`,
		dto.ID,
		projectID,
		taskID,
		dto.OperationType,
		dto.UserId,
		dto.CreatedAt,
	)
	return err
}

func (repository *historyOperationRepository) searchHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, scope scope, scopeID string) ([]domain.HistoryOperation, error) {
	var sqlQuery string
	switch scope {
	case scopeProject:
		sqlQuery = `
				SELECT
					HO.id, HO.operation_user_id, U.name, HO.operation_date, HO.operation_type
				FROM history_operations HO
				INNER JOIN users U ON U.id = HO.operation_user_id
				WHERE HO.project_id = ?
				ORDER BY HO.operation_date DESC
			`
	case scopeTask:
		sqlQuery = `
				SELECT
					HO.id, HO.operation_user_id, U.name, HO.operation_date, HO.operation_type
				FROM history_operations HO
				INNER JOIN users U ON U.id = HO.operation_user_id
				WHERE HO.task_id = ?
				ORDER BY HO.operation_date DESC
			`
	default:
		return nil, fmt.Errorf("invalid history operation scope: %s", scope)
	}
	rows, err := dbExecutor.QueryContext(ctx, sqlQuery, scopeID)
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

func (repository *historyOperationRepository) AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, operation domain.HistoryOperation) error {
	return repository.addHistoryOperation(ctx, dbExecutor, projectID, nil, operation)
}

func (repository *historyOperationRepository) SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectID string) ([]domain.HistoryOperation, error) {
	return repository.searchHistoryOperations(ctx, dbExecutor, scopeProject, projectID)
}

func (repository *historyOperationRepository) AddTaskHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectID string, taskID string, operation domain.HistoryOperation) error {
	return repository.addHistoryOperation(ctx, dbExecutor, projectID, &taskID, operation)
}

func (repository *historyOperationRepository) SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskID string) ([]domain.HistoryOperation, error) {
	return repository.searchHistoryOperations(ctx, dbExecutor, scopeTask, taskID)
}
