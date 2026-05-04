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
	GetWorkspaceProjectTypes(ctx context.Context, workspaceId string) ([]projectTypeDTO, error)
	Search(ctx context.Context, filter projectTypeFilterDTO) ([]projectTypeDTO, error)
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
            INSERT INTO project_types (id, name, item_hex_color)
			VALUES (?, ?, ?, ?)
        `,
		projectType.ID,
		projectType.Name,
		projectType.HexColor,
	)
	return err
}

func (projectTypeRepository *projectTypeRepository) Update(ctx context.Context, projectType projectTypeDTO) error {
	_, err := projectTypeRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_types SET
				name = ?,
				item_hex_color = ?
			WHERE id = ?
        `,
		projectType.Name,
		projectType.HexColor,
		projectType.ID,
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
                PT.id, PT.name, PT.item_hex_color
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

func (projectTypeRepository *projectTypeRepository) Search(ctx context.Context, filter projectTypeFilterDTO) ([]projectTypeDTO, error) {
	rows, err := projectTypeRepository.database.QueryContext(
		ctx,
		`
			SELECT
				PT.id, PT.name, PT.item_hex_color
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
			&projectType.ID, &projectType.Name, &projectType.HexColor,
		); err != nil {
			return nil, err
		}
		projectTypes = append(projectTypes, projectType)
	}
	return projectTypes, nil
}
