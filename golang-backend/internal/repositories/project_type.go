package repositories

import (
	"context"
	"database/sql"

	"github.com/aportela/gotask/internal/models"
)

type ProjectTypeRepository struct {
	database *sql.DB
}

func NewProjectTypeRepository(database *sql.DB) *ProjectTypeRepository {
	return &ProjectTypeRepository{
		database: database,
	}
}

func (projectTypeRepository *ProjectTypeRepository) Add(ctx context.Context, projectType models.ProjectType) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO PROJECT_TYPE (id, name) VALUES (?, ?)
        `,
		projectType.ID,
		projectType.Name,
	)
	return err
}

func (projectTypeRepository *ProjectTypeRepository) Update(ctx context.Context, projectType models.ProjectType) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            UPDATE PROJECT_TYPE SET name = ? WHERE id = ?
        `,
		projectType.ID,
		projectType.Name,
	)
	return err
}

func (projectTypeRepository *ProjectTypeRepository) Delete(ctx context.Context, id string) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM PROJECT_TYPE
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (ProjectTypeRepository *ProjectTypeRepository) Search(ctx context.Context) ([]models.ProjectType, error) {
	rows, err := ProjectTypeRepository.database.QueryContext(
		ctx,
		`
			SELECT
					PT.id, PT.name
			FROM PROJECT_TYPE PT
			ORDER BY PT.name
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectTypes []models.ProjectType

	for rows.Next() {
		var projectType models.ProjectType

		if err := rows.Scan(
			&projectType.ID, &projectType.Name,
		); err != nil {
			return nil, err
		}

		projectTypes = append(projectTypes, projectType)
	}

	return projectTypes, nil
}
