package taskrelationrepository

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type TaskRelationRepository interface {
	AddTaskRelation(ctx context.Context, dbExecutor database.DatabaseExecutor, taskRelation domain.TaskRelation) error
	DeleteTaskRelation(ctx context.Context, dbExecutor database.DatabaseExecutor, taskRelation domain.TaskRelation) error
}

type taskRelationRepository struct {
}

func NewRepository() TaskRelationRepository {
	return &taskRelationRepository{}
}

func (repository *taskRelationRepository) AddTaskRelation(ctx context.Context, dbExecutor database.DatabaseExecutor, taskRelation domain.TaskRelation) error {
	dto := toDTO(taskRelation)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO task_relations
				(task_id, related_task_id, relation_type )
			VALUES
				(?, ?)
        `,
		dto.TaskID,
		dto.RelatedTaskID,
		dto.RelationType,
	)
	return err
}

func (repository *taskRelationRepository) DeleteTaskRelation(ctx context.Context, dbExecutor database.DatabaseExecutor, taskRelation domain.TaskRelation) error {
	var sqlQuery string
	var queryArgs []any
	switch taskRelation.RelationType {
	case domain.RelationTypeLink:
		sqlQuery = `
            DELETE FROM task_relations
			WHERE (
					task_id = ?
				AND
					related_task_id = ?
				AND
					relation_type = ?
			)
			OR (
					task_id = ?
				AND
					related_task_id = ?
				AND
					relation_type = ?
			)
        `
		queryArgs = append(queryArgs, taskRelation.TaskID, taskRelation.RelatedTaskID, taskRelation.RelationType, taskRelation.RelatedTaskID, taskRelation.TaskID, taskRelation.RelationType)
	case domain.RelationTypeChild:
		sqlQuery = `
            DELETE FROM task_relations
			WHERE (
					task_id = ?
				AND
					related_task_id = ?
				AND
					relation_type = ?
			)

        `
		queryArgs = append(queryArgs, taskRelation.TaskID, taskRelation.RelatedTaskID, taskRelation.RelationType)
	default:
		return fmt.Errorf("Invalid task relation type")
	}
	_, err := dbExecutor.ExecContext(ctx, sqlQuery, queryArgs)
	return err
}
