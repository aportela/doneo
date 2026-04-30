package projectrepository

import (
	"context"
	"database/sql"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectRepository interface {
	Add(ctx context.Context, project projectDTO) error
	Update(ctx context.Context, project projectDTO) error
	Get(ctx context.Context, id string) (projectDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]projectDTO, error)
}

type projectRepository struct {
	database database.Database
}

func NewProjectRepository(database database.Database) ProjectRepository {
	return &projectRepository{database: database}
}

func (projectRepository *projectRepository) Add(ctx context.Context, project projectDTO) error {
	_, err := projectRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO projects
				(id, workspace_id, key, summary, description, creator_id, created_at, updated_at, started_at, finished_at, due_at, priority_id, status_id, type_id)
			VALUES
				(?, ?, ?, ?, ?, ?, ?, NULL, ?, ?, ?, ?, ?, ?)
        `,
		project.ID,
		project.WorkspaceId,
		project.Key,
		project.Summary,
		utils.NullableStringToSQL(project.Description),
		project.CreatorId,
		project.CreatedAt,
		utils.NullableInt64ToSQL(project.StartedAt),
		utils.NullableInt64ToSQL(project.FinishedAt),
		utils.NullableInt64ToSQL(project.DueAt),
		project.PriorityId,
		project.StatusId,
		project.TypeId,
	)
	return err
}

func (projectRepository *projectRepository) Update(ctx context.Context, project projectDTO) error {
	_, err := projectRepository.database.ExecContext(
		ctx,
		`
            UPDATE projects SET
				key = ?,
				summary = ?,
				description = ?,
				updated_at = ?,
				started_at = ?,
				finished_at = ?,
				due_at = ?,
				priority_id = ?,
				status_id = ?,
				type_id = ?
			WHERE id = ?
        `,
		project.Key,
		project.Summary,
		utils.NullableStringToSQL(project.Description),
		utils.CurrentMSTimestamp(),
		utils.NullableInt64ToSQL(project.StartedAt),
		utils.NullableInt64ToSQL(project.UpdatedAt),
		utils.NullableInt64ToSQL(project.DueAt),
		project.PriorityId,
		project.StatusId,
		project.TypeId,
		project.ID,
	)
	return err
}

func (projectRepository *projectRepository) Delete(ctx context.Context, id string) error {
	_, err := projectRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM projects
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (projectRepository *projectRepository) Get(ctx context.Context, id string) (projectDTO, error) {
	var project projectDTO
	var description sql.NullString
	var updated_at, started_at, finished_at, due_at sql.NullInt64
	err := projectRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.created_at, P.updated_at, P.started_at, P.finished_at, P.due_at, P.status_id, PS.name AS status_name, P.priority_id, PP.name AS priority_name, P.type_id, PT.name AS type_name, P.creator_id, U.name AS creator_name
            FROM projects P
			INNER JOIN project_priorities PP ON PP.id = P.priority_id
			INNER JOIN project_statuses PS ON PS.id = P.status_id
			INNER JOIN project_types PT ON PT.id = P.type_id
			INNER JOIN users U ON U.ID = P.creator_id
            WHERE P.id = ?
        `,
		id).Scan(
		&project.ID, &project.Key, &project.Summary, &description,
		&project.CreatedAt, &updated_at, &started_at, &finished_at, &due_at,
		&project.StatusId, &project.StatusName, &project.PriorityId, &project.PriorityName,
		&project.TypeId, &project.TypeName, &project.CreatorId, &project.CreatorName,
	)
	project.Description = utils.SQLStrPtr(description)
	project.UpdatedAt = utils.SQLInt64Ptr(updated_at)
	project.StartedAt = utils.SQLInt64Ptr(started_at)
	project.UpdatedAt = utils.SQLInt64Ptr(finished_at)
	project.DueAt = utils.SQLInt64Ptr(due_at)
	return project, err
}

func (projectRepository *projectRepository) Search(ctx context.Context) ([]projectDTO, error) {
	rows, err := projectRepository.database.QueryContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.created_at, P.updated_at, P.started_at, P.finished_at, P.due_at, P.status_id, PS.name AS status_name, P.priority_id, PP.name AS priority_name, P.type_id, PT.name AS type_name, P.creator_id, U.name AS creator_name
            FROM projects P
			INNER JOIN project_priorities PP ON PP.id = P.priority_id
			INNER JOIN project_statuses PS ON PS.id = P.status_id
			INNER JOIN project_types PT ON PT.id = P.type_id
			INNER JOIN users U ON U.ID = P.creator_id
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projects []projectDTO
	for rows.Next() {
		var project projectDTO
		var description sql.NullString
		var updated_at, started_at, finished_at, due_at sql.NullInt64
		if err := rows.Scan(
			&project.ID, &project.Key, &project.Summary, &description,
			&project.CreatedAt, &updated_at, &started_at, &finished_at, &due_at,
			&project.StatusId, &project.StatusName, &project.PriorityId, &project.PriorityName,
			&project.TypeId, &project.TypeName, &project.CreatorId, &project.CreatorName,
		); err != nil {
			return nil, err
		}
		project.Description = utils.SQLStrPtr(description)
		project.UpdatedAt = utils.SQLInt64Ptr(updated_at)
		project.StartedAt = utils.SQLInt64Ptr(started_at)
		project.UpdatedAt = utils.SQLInt64Ptr(finished_at)
		project.DueAt = utils.SQLInt64Ptr(due_at)
		projects = append(projects, project)
	}
	return projects, nil
}
