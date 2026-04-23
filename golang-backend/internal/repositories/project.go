package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/utils"
)

type ProjectRepository struct {
	database *sql.DB
}

func NewProjectRepository(database *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		database: database,
	}
}

func (projectRepository *ProjectRepository) Get(ctx context.Context, id string) (models.Project, error) {
	var project models.Project
	var lmtime, stime, ftime, dtime sql.NullInt64
	var description sql.NullString
	err := projectRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.ctime, P.lmtime, P.stime, P.ftime, P.dtime, P.type, PT.name
            FROM PROJECT P
			LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
            WHERE P.id = ?
        `,
		id).Scan(&project.ID, &project.Key, &project.Summary, &description, &project.CreatedAt, &lmtime, &stime, &ftime, &dtime, &project.Type.ID, &project.Type.Name)
	project.Description = utils.StrPtr(description)
	project.LastModifiedAt = utils.Int64Ptr(lmtime)
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
                P.id, P.key, P.summary, P.description, P.ctime, P.lmtime, P.stime, P.ftime, P.dtime, P.type, PT.name
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
		var lmtime, stime, ftime, dtime sql.NullInt64
		var description sql.NullString

		if err := rows.Scan(
			&project.ID, &project.Key, &project.Summary, &description,
			&project.CreatedAt, &lmtime, &stime, &ftime, &dtime,
			&project.Type.ID, &project.Type.Name,
		); err != nil {
			return nil, err
		}

		project.Description = utils.StrPtr(description)
		project.LastModifiedAt = utils.Int64Ptr(lmtime)
		project.StartedAt = utils.Int64Ptr(stime)
		project.FinishedAt = utils.Int64Ptr(ftime)
		project.DueAt = utils.Int64Ptr(dtime)

		projects = append(projects, project)
	}

	return projects, nil
}
