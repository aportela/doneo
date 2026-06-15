package historyoperationrepository

import (
	"context"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories"
)

type HistoryOperationRepository interface {
	AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, operation domain.HistoryOperation) error
	SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectId string) ([]domain.HistoryOperation, error)
	AddTaskOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, taskId string, operation domain.HistoryOperation) error
	SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskId string) ([]domain.HistoryOperation, error)
}

type historyOperationRepository struct{}

func NewRepository() HistoryOperationRepository {
	return &historyOperationRepository{}
}

func (repository *historyOperationRepository) AddProjectHistoryOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, operation domain.HistoryOperation) error {
	dto := toDTO(operation)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
			INSERT INTO history_operations
				(id, project_id, task_id, operation_type, operation_user_id, operation_date)
			VALUES
				(?, ?, NULL, ?, ?, ?)
		`,
		dto.ID,
		projectId,
		dto.OperationType,
		dto.UserId,
		dto.CreatedAt,
	)
	return err
}

func (repository *historyOperationRepository) SearchProjectHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, projectId string) ([]domain.HistoryOperation, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				HO.id, HO.operation_user_id, U.name, HO.operation_date, HO.operation_type
            FROM history_operations HO
			INNER JOIN users U ON U.id = HO.operation_user_id
            WHERE HO.project_id = ?
			ORDER BY HO.operation_date DESC
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

func (repository *historyOperationRepository) AddTaskOperation(ctx context.Context, dbExecutor repositories.Executor, projectId string, taskId string, operation domain.HistoryOperation) error {
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
		projectId,
		taskId,
		dto.OperationType,
		dto.UserId,
		dto.CreatedAt,
	)
	return err
}

func (repository *historyOperationRepository) SearchTaskHistoryOperations(ctx context.Context, dbExecutor repositories.Executor, taskId string) ([]domain.HistoryOperation, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
				HO.id, HO.operation_user_id, U.name, HO.operation_date, HO.operation_type
            FROM history_operations HO
			INNER JOIN users U ON U.id = HO.operation_user_id
            WHERE HO.task_id = ?
			ORDER BY HO.operation_date DESC
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
