package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/utils"
)

type ProjectRepository struct {
	database database.Database
}

func NewProjectRepository(database database.Database) *ProjectRepository {
	return &ProjectRepository{
		database: database,
	}
}

func (projectRepository *ProjectRepository) Add(ctx context.Context, project models.Project) error {
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
		utils.CurrentTimestamp(),
		utils.NullableInt64ToSQL(project.LastModifiedAt),
		utils.NullableInt64ToSQL(project.StartedAt),
		utils.NullableInt64ToSQL(project.FinishedAt),
		utils.NullableInt64ToSQL(project.DueAt),
		project.Type.ID,
	)
	return err
}

func (projectRepository *ProjectRepository) Update(ctx context.Context, project models.Project) error {
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
		utils.CurrentTimestamp(),
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

func (projectRepository *ProjectRepository) Get(ctx context.Context, id string) (models.Project, error) {
	var project models.Project
	var mtime, stime, ftime, dtime sql.NullInt64
	var description sql.NullString
	err := projectRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.ctime, P.mtime, P.stime, P.ftime, P.dtime, P.type, PT.name
            FROM PROJECT P
			LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
            WHERE P.id = ?
        `,
		id).Scan(&project.ID, &project.Key, &project.Summary, &description, &project.CreatedAt, &mtime, &stime, &ftime, &dtime, &project.Type.ID, &project.Type.Name)
	project.Description = utils.StrPtr(description)
	project.LastModifiedAt = utils.Int64Ptr(mtime)
	project.StartedAt = utils.Int64Ptr(stime)
	project.FinishedAt = utils.Int64Ptr(ftime)
	project.DueAt = utils.Int64Ptr(dtime)

	return project, err
}

func (projectRepository *ProjectRepository) Search(ctx context.Context) ([]models.Project, error) {
	rows, err := projectRepository.database.QueryContext(
		ctx,
		`
        SELECT
                P.id, P.key, P.summary, P.description, P.ctime, P.mtime, P.stime, P.ftime, P.dtime, P.type, PT.name
		FROM PROJECT P
		LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var project models.Project
		var mtime, stime, ftime, dtime sql.NullInt64
		var description sql.NullString

		if err := rows.Scan(
			&project.ID, &project.Key, &project.Summary, &description,
			&project.CreatedAt, &mtime, &stime, &ftime, &dtime,
			&project.Type.ID, &project.Type.Name,
		); err != nil {
			return nil, err
		}

		project.Description = utils.StrPtr(description)
		project.LastModifiedAt = utils.Int64Ptr(mtime)
		project.StartedAt = utils.Int64Ptr(stime)
		project.FinishedAt = utils.Int64Ptr(ftime)
		project.DueAt = utils.Int64Ptr(dtime)

		projects = append(projects, project)
	}

	return projects, nil
}
