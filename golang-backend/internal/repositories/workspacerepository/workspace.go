package workspacerepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
)

type WorkspaceRepository interface {
	Add(ctx context.Context, workspace workspaceDTO) error
	Update(ctx context.Context, workspace workspaceDTO) error
	Get(ctx context.Context, id string) (workspaceDTO, error)
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context) ([]workspaceDTO, error)
}

type workspaceRepository struct {
	database database.Database
}

func NewWorkspaceRepository(database database.Database) WorkspaceRepository {
	return &workspaceRepository{database: database}
}

func (workspaceRepository *workspaceRepository) Add(ctx context.Context, workspace workspaceDTO) error {
	_, err := workspaceRepository.database.ExecContext(
		ctx,
		`
            INSERT INTO workspaces (id, name, description, creator_id, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, NULL)
        `,
		workspace.ID,
		workspace.Name,
		utils.NullableStringToSQL(workspace.Description),
		workspace.CreatorId,
		workspace.CreatedAt,
	)
	return err
}

func (workspaceRepository *workspaceRepository) Update(ctx context.Context, workspace workspaceDTO) error {
	_, err := workspaceRepository.database.ExecContext(
		ctx,
		`
            UPDATE workspaces SET
				name = ?,
				description = ?
				updated_at = ?
			WHERE id = ?
        `,
		workspace.Name,
		utils.NullableStringToSQL(workspace.Description),
		utils.NullableInt64ToSQL(workspace.UpdatedAt),
		workspace.ID,
	)
	return err
}

func (workspaceRepository *workspaceRepository) Delete(ctx context.Context, id string) error {
	_, err := workspaceRepository.database.ExecContext(
		ctx,
		`
            DELETE FROM workspaces
			WHERE id = ?
        `,
		id,
	)
	return err
}

func (workspaceRepository *workspaceRepository) Get(ctx context.Context, id string) (workspaceDTO, error) {
	var workspace workspaceDTO
	var description sql.NullString
	var updatedAt sql.NullInt64
	err := workspaceRepository.database.QueryRowContext(
		ctx,
		`
            SELECT
                W.id, W.name, W.description, W.created_at, W.updated_at, W.creator_id, U.name
            FROM workspaces W
			INNER JOIN users U ON U.id = W.creator_id
            WHERE W.id = ?
        `,
		id).Scan(&workspace.ID, &workspace.Name, &description, &workspace.CreatedAt, &updatedAt, &workspace.CreatorId, &workspace.CreatorName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workspace, domain.ErrNotFound
		}
		return workspace, err
	}
	if description.Valid {
		workspace.Description = &description.String
	}
	workspace.UpdatedAt = utils.SQLInt64Ptr(updatedAt)
	return workspace, err
}

func (workspaceRepository *workspaceRepository) Search(ctx context.Context) ([]workspaceDTO, error) {
	rows, err := workspaceRepository.database.QueryContext(
		ctx,
		`
			SELECT
				W.id, W.name, W.description, W.creator_id, U.name AS creator_name, W.created_at, W.updated_at
			FROM workspaces W
			INNER JOIN users U ON U.id = W.creator_id
			ORDER BY W.name
        `,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var workspaces []workspaceDTO
	for rows.Next() {
		var workspace workspaceDTO
		var description sql.NullString
		var updatedAt sql.NullInt64
		if err := rows.Scan(&workspace.ID, &workspace.Name, &description, &workspace.CreatorId, &workspace.CreatorName, &workspace.CreatedAt, &updatedAt); err != nil {
			return nil, err
		}
		if description.Valid {
			workspace.Description = &description.String
		}
		workspace.UpdatedAt = utils.SQLInt64Ptr(updatedAt)
		workspaces = append(workspaces, workspace)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return workspaces, nil
}
