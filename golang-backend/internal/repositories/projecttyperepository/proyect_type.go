package projecttyperepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProjectTypeRepository interface {
	Add(ctx context.Context, projectType projectTypeDTO) error
	Update(ctx context.Context, projectType projectTypeDTO) error
	Get(ctx context.Context, id string) (projectTypeDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]projectTypeDTO, error)
}

type projectTypeRepository struct {
	database database.Database
}

func NewProjectTypeRepository(database database.Database) ProjectTypeRepository {
	return &projectTypeRepository{database: database}
}

func (projectTypeRepository *projectTypeRepository) Add(ctx context.Context, projectType projectTypeDTO) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_types (id, name)
			VALUES (?, ?)
        `,
		projectType.ID,
		projectType.Name,
	)
	return err
}

func (projectTypeRepository *projectTypeRepository) Update(ctx context.Context, projectType projectTypeDTO) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_types SET
				name = ?
			WHERE id = ?
        `,
		projectType.ID,
		projectType.Name,
	)
	return err
}

func (projectTypeRepository *projectTypeRepository) Delete(ctx context.Context, id string) error {
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

func (projectTypeRepository *projectTypeRepository) Get(ctx context.Context, id string) (projectTypeDTO, error) {
	var projectType projectTypeDTO
	err := projectTypeRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PT.id, PT.name
            FROM project_types PT
            WHERE PT.id = ?
        `,
		id).Scan(&projectType.ID, &projectType.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectType, domain.ErrNotFound
		}
		return projectType, err
	}
	return projectType, err
}

func (projectTypeRepository *projectTypeRepository) Search(ctx context.Context) ([]projectTypeDTO, error) {
	rows, err := projectTypeRepository.database.QueryContext(
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
	var projectTypes []projectTypeDTO
	for rows.Next() {
		var projectType projectTypeDTO

		if err := rows.Scan(
			&projectType.ID, &projectType.Name,
		); err != nil {
			return nil, err
		}

		projectTypes = append(projectTypes, projectType)
	}
	return projectTypes, nil
}
