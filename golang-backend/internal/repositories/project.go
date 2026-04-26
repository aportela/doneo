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
            INSERT INTO PROJECT_PARTICIPANT (project_id, user_id)
			VALUES (?, ?)
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
            INSERT INTO PROJECT (id, key, summary, description, cuser, ctime, mtime, stime, ftime, dtime, type)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        `,
		project.ID,
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
            UPDATE PROJECT SET
				key = ?,
				summary = ?,
				description = ?,
				mtime = ?,
				stime = ?,
				ftime = ?,
				dtime = ?,
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
            DELETE FROM PROJECT
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
		FROM PROJECT_PARTICIPANT PP
		INNER JOIN USER U ON U.id = PP.user_id
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
	var mtime, stime, ftime, dtime sql.NullInt64
	var description sql.NullString
	var creatorID, creatorName string
	err := projectRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.ctime, P.mtime, P.stime, P.ftime, P.dtime, P.type, PT.name, P.cuser, U.name
            FROM PROJECT P
			LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
			INNER JOIN USER U ON U.ID = P.cuser
            WHERE P.id = ?
        `,
		id).Scan(&project.ID, &project.Key, &project.Summary, &description, &project.CreatedAt, &mtime, &stime, &ftime, &dtime, &project.Type.ID, &project.Type.Name, &creatorID, &creatorName)
	project.CreatedBy = domain.UserBase{ID: creatorID, Name: creatorName}
	project.Description = utils.StrPtr(description)
	project.LastModifiedAt = utils.Int64Ptr(mtime)
	project.StartedAt = utils.Int64Ptr(stime)
	project.FinishedAt = utils.Int64Ptr(ftime)
	project.DueAt = utils.Int64Ptr(dtime)

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
                P.id, P.key, P.summary, P.description, P.ctime, P.mtime, P.stime, P.ftime, P.dtime, P.type, PT.name, P.cuser, U.name
		FROM PROJECT P
		LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
		INNER JOIN USER U ON U.ID = P.cuser
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []domain.Project

	for rows.Next() {
		var project domain.Project
		var mtime, stime, ftime, dtime sql.NullInt64
		var description sql.NullString
		var creatorID, creatorName string

		if err := rows.Scan(
			&project.ID, &project.Key, &project.Summary, &description,
			&project.CreatedAt, &mtime, &stime, &ftime, &dtime,
			&project.Type.ID, &project.Type.Name, &creatorID, &creatorName,
		); err != nil {
			return nil, err
		}

		project.CreatedBy = domain.UserBase{ID: creatorID, Name: creatorName}
		project.Description = utils.StrPtr(description)
		project.LastModifiedAt = utils.Int64Ptr(mtime)
		project.StartedAt = utils.Int64Ptr(stime)
		project.FinishedAt = utils.Int64Ptr(ftime)
		project.DueAt = utils.Int64Ptr(dtime)

		projects = append(projects, project)
	}

	return projects, nil
}
