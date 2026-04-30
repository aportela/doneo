package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

type ProjectRepository struct {
	database database.Database
}

func NewProjectRepository(database database.Database) *ProjectRepository {
	return &ProjectRepository{
		database: database,
	}
}

func (projectRepository *ProjectRepository) AddParticipant(ctx context.Context, projectId string, userId string) error {
	_, err := projectRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_user_roles (project_id, user_id, is_admin, is_member, is_guest)
			VALUES (?, ?, 1, 1, 1)
        `,
		projectId,
		userId,
	)
	return err
}

func (projectRepository *ProjectRepository) add(ctx context.Context, project domain.Project) error {
	_, err := projectRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO projects (id, workspace_id, key, summary, description, creator_id, created_at, updated_at, started_at, finished_at, due_at, type)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `,
		project.ID,
		"019dddbc-9df2-717c-be35-70f707e6098f",
		project.Key,
		project.Summary,
		utils.NullableStringToSQL(project.Description),
		project.CreatedBy.ID,
		project.CreatedAt,
		utils.NullableInt64ToSQL(project.LastModifiedAt),
		utils.NullableInt64ToSQL(project.StartedAt),
		utils.NullableInt64ToSQL(project.FinishedAt),
		utils.NullableInt64ToSQL(project.DueAt),
		project.Type.ID,
	)
	return err
}

func (projectRepository *ProjectRepository) Add(ctx context.Context, project domain.Project) error {
	// TODO: transaction
	err := projectRepository.add(ctx, project)
	if err != nil {
		return err
	}
	for _, projectParticipant := range project.Participants {
		err = projectRepository.AddParticipant(ctx, project.ID, projectParticipant.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (projectRepository *ProjectRepository) Update(ctx context.Context, project domain.Project) error {
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
				type = ?
			WHERE id = ?
        `,
		project.Key,
		project.Summary,
		utils.NullableStringToSQL(project.Description),
		utils.CurrentMSTimestamp(),
		utils.NullableInt64ToSQL(project.StartedAt),
		utils.NullableInt64ToSQL(project.FinishedAt),
		utils.NullableInt64ToSQL(project.DueAt),
		project.Type.ID,
		project.ID,
	)
	return err
}

func (projectRepository *ProjectRepository) Delete(ctx context.Context, id string) error {
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

func (projectRepository *ProjectRepository) getParticipants(ctx context.Context, projectId string) ([]domain.UserBase, error) {
	rows, err := projectRepository.database.QueryContext(
		ctx,
		`
        SELECT
                U.id, U.name
		FROM project_user_roles PP
		INNER JOIN users U ON U.id = PP.user_id
		WHERE PP.project_id = ?
		ORDER BY U.name
        `,
		projectId,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []domain.UserBase

	for rows.Next() {
		var user domain.UserBase

		if err := rows.Scan(
			&user.ID, &user.Name,
		); err != nil {
			return nil, err
		}

		participants = append(participants, user)
	}

	return participants, nil
}

func (projectRepository *ProjectRepository) get(ctx context.Context, id string) (*domain.Project, error) {
	var project domain.Project
	var updated_at, started_at, finished_at, due_at sql.NullInt64
	var description sql.NullString
	var creatorID, creatorName string
	err := projectRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.created_at, P.updated_at, P.started_at, P.finished_at, P.due_at, P.type, PT.name, P.creator_id, U.name
            FROM projects P
			LEFT JOIN project_types PT ON PT.id = P.type
			INNER JOIN users U ON U.ID = P.creator_id
            WHERE P.id = ?
        `,
		id).Scan(&project.ID, &project.Key, &project.Summary, &description, &project.CreatedAt, &updated_at, &started_at, &finished_at, &due_at, &project.Type.ID, &project.Type.Name, &creatorID, &creatorName)
	project.CreatedBy = domain.UserBase{ID: creatorID, Name: creatorName}
	project.Description = utils.SQLStrPtr(description)
	project.LastModifiedAt = utils.SQLInt64Ptr(updated_at)
	project.StartedAt = utils.SQLInt64Ptr(started_at)
	project.FinishedAt = utils.SQLInt64Ptr(finished_at)
	project.DueAt = utils.SQLInt64Ptr(due_at)

	return &project, err
}

func (projectRepository *ProjectRepository) Get(ctx context.Context, id string) (*domain.Project, error) {
	project, err := projectRepository.get(ctx, id)
	if err != nil {
		return nil, err
	}
	var participants []domain.UserBase
	participants, err = projectRepository.getParticipants(ctx, project.ID)
	if err != nil {
		return nil, err
	}
	project.Participants = participants
	return project, nil
}

func (projectRepository *ProjectRepository) Search(ctx context.Context) ([]domain.Project, error) {
	rows, err := projectRepository.database.QueryContext(
		ctx,
		`
        SELECT
                P.id, P.key, P.summary, P.description, P.created_at, P.updated_at, P.started_at, P.finished_at, P.due_at, P.type, PT.name, P.creator_id, U.name
		FROM projects P
		LEFT JOIN project_types PT ON PT.id = P.type
		INNER JOIN users U ON U.ID = P.creator_id
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project

	for rows.Next() {
		var project domain.Project
		var updated_at, started_at, finished_at, due_at sql.NullInt64
		var description sql.NullString
		var creatorID, creatorName string

		if err := rows.Scan(
			&project.ID, &project.Key, &project.Summary, &description,
			&project.CreatedAt, &updated_at, &started_at, &finished_at, &due_at,
			&project.Type.ID, &project.Type.Name, &creatorID, &creatorName,
		); err != nil {
			return nil, err
		}

		project.CreatedBy = domain.UserBase{ID: creatorID, Name: creatorName}
		project.Description = utils.SQLStrPtr(description)
		project.LastModifiedAt = utils.SQLInt64Ptr(updated_at)
		project.StartedAt = utils.SQLInt64Ptr(started_at)
		project.FinishedAt = utils.SQLInt64Ptr(finished_at)
		project.DueAt = utils.SQLInt64Ptr(due_at)

		projects = append(projects, project)
	}

	return projects, nil
}
