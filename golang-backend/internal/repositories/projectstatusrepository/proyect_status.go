package projectstatusrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProjectStatusRepository interface {
	Add(ctx context.Context, projectStatus projectStatusDTO) error
	Update(ctx context.Context, projectStatus projectStatusDTO) error
	Get(ctx context.Context, id string) (projectStatusDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]projectStatusDTO, error)
}

type projectStatusRepository struct {
	database database.Database
}

func NewProjectStatusRepository(database database.Database) ProjectStatusRepository {
	return &projectStatusRepository{database: database}
}

func (projectStatusRepository *projectStatusRepository) Add(ctx context.Context, projectStatus projectStatusDTO) error {
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_statuses (id, name, item_index, item_hex_color)
			VALUES (?, ?, ?, ?)
        `,
		projectStatus.ID,
		projectStatus.Name,
		projectStatus.Index,
		projectStatus.HexColor,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Update(ctx context.Context, projectStatus projectStatusDTO) error {
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_statuses SET
				name = ?,
				item_index = ?,
				item_hex_color = ?,
			WHERE id = ?
        `,
		projectStatus.Name,
		projectStatus.Index,
		projectStatus.HexColor,
		projectStatus.ID,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Delete(ctx context.Context, id string) error {
	_, err := projectStatusRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_statuses
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (projectStatusRepository *projectStatusRepository) Get(ctx context.Context, id string) (projectStatusDTO, error) {
	var projectStatus projectStatusDTO
	err := projectStatusRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PS.id, PS.name, PS.item_index, PS.item_hex_color
            FROM project_statuses PS
            WHERE PS.id = ?
        `,
		id).Scan(&projectStatus.ID, &projectStatus.Name, &projectStatus.Index, &projectStatus.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectStatus, domain.ErrNotFound
		}
		return projectStatus, err
	}
	return projectStatus, err
}

func (projectStatusRepository *projectStatusRepository) Search(ctx context.Context) ([]projectStatusDTO, error) {
	rows, err := projectStatusRepository.database.QueryContext(
		ctx,
		`
			SELECT
				PS.id, PS.name, PS.item_index, PS.item_hex_color
			FROM project_statuses PS
			ORDER BY PS.item_index, PS.name
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var projectStatuses []projectStatusDTO
	for rows.Next() {
		var projectStatus projectStatusDTO
		if err := rows.Scan(
			&projectStatus.ID, &projectStatus.Name, &projectStatus.Index, &projectStatus.HexColor,
		); err != nil {
			return nil, err
		}
		projectStatuses = append(projectStatuses, projectStatus)
	}
	return projectStatuses, nil
}
