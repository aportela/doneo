package repositories

import (
	"database/sql"

	"github.com/aportela/gotask/internal/models"
)

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

func (r *ProjectRepository) GetByID(id string) (models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(
		`
            SELECT
                P.id, P.key, P.summary, P.description, P.ctime, P.lmtime, P.stime, P.ftime, P.dtime, P.type, PT.name
            FROM PROJECT P
			LEFT JOIN PROJECT_TYPE PT ON PT.id = P.type
            WHERE P.id = ?
        `,
		id).Scan(&p.ID, &p.Key, &p.Summary, &p.Description, &p.CreatedAt, &p.LastModifiedAt, &p.StartedAt, &p.FinishedAt, &p.DueAt, &p.Type.ID, &p.Type.Name)

	return p, err
}

func (r *ProjectRepository) Search() ([]models.Project, error) {
	rows, err := r.db.Query(
		`
        SELECT
            id, summary
        FROM PROJECT
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project

	for rows.Next() {
		var t models.Project
		if err := rows.Scan(&t.ID, &t.Summary); err != nil {
			return nil, err
		}
		projects = append(projects, t)
	}

	return projects, nil
}
