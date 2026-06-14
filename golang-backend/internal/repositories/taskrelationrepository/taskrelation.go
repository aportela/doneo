package taskrelationrepository

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

type TaskRelationRepository interface {
	AddTaskRelation(ctx context.Context, taskRelation domain.TaskRelation) error
	DeleteTaskRelation(ctx context.Context, taskRelation domain.TaskRelation) error
}

type taskRelationRepository struct {
	db database.Database
}

func NewRepository(db database.Database) TaskRelationRepository {
	return &taskRelationRepository{db: db}
}

func (repository *taskRelationRepository) AddTaskRelation(ctx context.Context, taskRelation domain.TaskRelation) error {
	dto := toDTO(taskRelation)
	_, err := repository.db.ExecContext(
		ctx,
		`
            INSERT INTO task_relations (task_id, related_task_id, relation_type )
			VALUES (?, ?)
        `,
		dto.TaskID,
		dto.RelatedTaskID,
		dto.RelationType,
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
			return &domain.ValidationError{Field: "task_id, related_task_id"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(task_id)") {
				return &domain.ValidationError{Field: "task_id"}
			} else if strings.Contains(sqlErr.Error(), "length(related_task_id)") {
				return &domain.ValidationError{Field: "related_task_id"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *taskRelationRepository) DeleteTaskRelation(ctx context.Context, taskRelation domain.TaskRelation) error {
	var sqlQuery string
	var queryArgs []any
	switch taskRelation.RelationType {
	case domain.RelationTypeLink:
		sqlQuery = `
            DELETE FROM task_relations
			WHERE
				(task_id = ? AND related_task_id = ? AND relation_type = ?)
			OR
				(task_id = ? AND related_task_id = ? AND relation_type = ?)
        `
		queryArgs = append(queryArgs, taskRelation.TaskID, taskRelation.RelatedTaskID, taskRelation.RelationType, taskRelation.RelatedTaskID, taskRelation.TaskID, taskRelation.RelationType)
	case domain.RelationTypeChild:
		sqlQuery = `
            DELETE FROM task_relations
			WHERE
				(task_id = ? AND related_task_id = ? AND relation_type = ?)

        `
		queryArgs = append(queryArgs, taskRelation.TaskID, taskRelation.RelatedTaskID, taskRelation.RelationType)
	default:
		// TODO:
		return fmt.Errorf("Invalid relation")
	}
	_, err := repository.db.ExecContext(
		ctx,
		sqlQuery,
		queryArgs,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}
