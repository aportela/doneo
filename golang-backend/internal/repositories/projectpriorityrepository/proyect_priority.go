package projectpriorityrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProjectPriorityRepository interface {
	Add(ctx context.Context, projectPriority projectPriorityDTO) error
	Update(ctx context.Context, projectPriority projectPriorityDTO) error
	Get(ctx context.Context, id string) (projectPriorityDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]projectPriorityDTO, error)
}

type projectPriorityRepository struct {
	database database.Database
}

func NewProjectPriorityRepository(database database.Database) ProjectPriorityRepository {
	return &projectPriorityRepository{database: database}
}

func (projectPriorityRepository *projectPriorityRepository) Add(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_priorities (id, name, item_index, item_hex_color)
			VALUES (?, ?, ?, ?)
        `,
		projectPriority.ID,
		projectPriority.Name,
		projectPriority.Index,
		projectPriority.HexColor,
	)
	return err
}

func (projectPriorityRepository *projectPriorityRepository) Update(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_priorities SET
				name = ?,
				item_index = ?,
				item_hex_color = ?
			WHERE id = ?
        `,
		projectPriority.Name,
		projectPriority.Index,
		projectPriority.HexColor,
		projectPriority.ID,
	)
	return err
}

func (projectPriorityRepository *projectPriorityRepository) Delete(ctx context.Context, id string) error {
	_, err := projectPriorityRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_priorities
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (projectPriorityRepository *projectPriorityRepository) Get(ctx context.Context, id string) (projectPriorityDTO, error) {
	var projectPriority projectPriorityDTO
	err := projectPriorityRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PP.id, PP.name, PP.item_index, PP.item_hex_color
            FROM project_priorities PP
            WHERE PP.id = ?
        `,
		id).Scan(&projectPriority.ID, &projectPriority.Name, &projectPriority.Index, &projectPriority.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return projectPriority, domain.ErrNotFound
		}
		return projectPriority, err
	}
	return projectPriority, err
}

func (projectPriorityRepository *projectPriorityRepository) Search(ctx context.Context) ([]projectPriorityDTO, error) {
	rows, err := projectPriorityRepository.database.QueryContext(
		ctx,
		`
			SELECT
				PP.id, PP.name, PP.item_index, PP.item_hex_color
			FROM project_priorities PP
			ORDER BY PP.item_index, PP.name
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var proyecPriorities []projectPriorityDTO
	for rows.Next() {
		var projectPriority projectPriorityDTO
		if err := rows.Scan(
			&projectPriority.ID, &projectPriority.Name, &projectPriority.Index, &projectPriority.HexColor,
		); err != nil {
			return nil, err
		}
		proyecPriorities = append(proyecPriorities, projectPriority)
	}
	return proyecPriorities, nil
}
