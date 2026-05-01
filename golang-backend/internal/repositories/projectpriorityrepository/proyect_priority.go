package projectpriorityrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProyectPriorityRepository interface {
	Add(ctx context.Context, projectPriority projectPriorityDTO) error
	Update(ctx context.Context, projectPriority projectPriorityDTO) error
	Get(ctx context.Context, id string) (projectPriorityDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]projectPriorityDTO, error)
}

type proyectPriorityRepository struct {
	database database.Database
}

func NewProyectPriorityRepository(database database.Database) ProyectPriorityRepository {
	return &proyectPriorityRepository{database: database}
}

func (proyectPriorityRepository *proyectPriorityRepository) Add(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := proyectPriorityRepository.database.ExecContext(
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

func (proyectPriorityRepository *proyectPriorityRepository) Update(ctx context.Context, projectPriority projectPriorityDTO) error {
	_, err := proyectPriorityRepository.database.ExecContext(
		ctx,
		`
            UPDATE project_priorities SET
				name = ?,
				item_index = ?,
				item_hex_color = ?
			WHERE id = ?
        `,
		projectPriority.ID,
		projectPriority.Name,
		projectPriority.Index,
		projectPriority.HexColor,
	)
	return err
}

func (proyectPriorityRepository *proyectPriorityRepository) Delete(ctx context.Context, id string) error {
	_, err := proyectPriorityRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_priorities
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (proyectPriorityRepository *proyectPriorityRepository) Get(ctx context.Context, id string) (projectPriorityDTO, error) {
	var proyectPriority projectPriorityDTO
	err := proyectPriorityRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PP.id, PP.name, PP.item_index, PP.item_hex_color
            FROM project_priorities PP
            WHERE PP.id = ?
        `,
		id).Scan(&proyectPriority.ID, &proyectPriority.Name, &proyectPriority.Index, &proyectPriority.HexColor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return proyectPriority, domain.ErrNotFound
		}
		return proyectPriority, err
	}
	return proyectPriority, err
}

func (proyectPriorityRepository *proyectPriorityRepository) Search(ctx context.Context) ([]projectPriorityDTO, error) {
	rows, err := proyectPriorityRepository.database.QueryContext(
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
		var proyectPriority projectPriorityDTO
		if err := rows.Scan(
			&proyectPriority.ID, &proyectPriority.Name, &proyectPriority.Index, &proyectPriority.HexColor,
		); err != nil {
			return nil, err
		}
		proyecPriorities = append(proyecPriorities, proyectPriority)
	}
	return proyecPriorities, nil
}
