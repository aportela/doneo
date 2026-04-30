package repositories

import (
	"context"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProjectTypeRepository struct {
	database database.Database
}

func NewProjectTypeRepository(database database.Database) *ProjectTypeRepository {
	return &ProjectTypeRepository{
		database: database,
	}
}

func (projectTypeRepository *ProjectTypeRepository) Add(ctx context.Context, projectType domain.ProjectType) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_types (id, name) VALUES (?, ?)
        `,
		projectType.ID,
		projectType.Name,
	)
	return err
}

func (projectTypeRepository *ProjectTypeRepository) Update(ctx context.Context, projectType domain.ProjectType) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_types SET name = ? WHERE id = ?
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
            DELETE FROM project_types
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (ProjectTypeRepository *ProjectTypeRepository) Search(ctx context.Context) ([]domain.ProjectType, error) {
	rows, err := ProjectTypeRepository.database.QueryContext(
		ctx,
		`
			SELECT
					PT.id, PT.name
			FROM project_types PT
			ORDER BY PT.name
        `,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projectTypes []domain.ProjectType

	for rows.Next() {
		var projectType domain.ProjectType

		if err := rows.Scan(
			&projectType.ID, &projectType.Name,
		); err != nil {
			return nil, err
		}

		projectTypes = append(projectTypes, projectType)
	}

	return projectTypes, nil
}
