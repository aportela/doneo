package projectpermissionrepository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
)

type ProjectPermissionRepository interface {
	Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, permission domain.ProjectPermission) error
	Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, permissionID string) error
	Get(ctx context.Context, dbExecutor database.DatabaseExecutor, permissionID string) (domain.ProjectPermission, error)
	GetProjectPermissions(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.ProjectPermission, error)
}

type projectPermissionRepository struct{}

func NewRepository() ProjectPermissionRepository {
	return &projectPermissionRepository{}
}

func (repository *projectPermissionRepository) Add(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string, permission domain.ProjectPermission) error {
	dto := toDTO(permission)
	_, err := dbExecutor.ExecContext(
		ctx,
		`
            INSERT INTO project_user_role
				(id, project_id, user_id, role_id)
			VALUES
				(?, ?, ?, ?)
        `,
		dto.ID,
		projectID,
		dto.UserID,
		dto.RoleID,
	)
	if err != nil {
		return mapSQLiteError(err)
	}
	return nil
}

func (repository *projectPermissionRepository) Delete(ctx context.Context, dbExecutor database.DatabaseExecutor, permissionID string) error {
	result, err := dbExecutor.ExecContext(
		ctx,
		`
            DELETE FROM project_user_role
			WHERE
				id = ?
        `,
		permissionID,
	)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count < 1 {
		return domain.NotFoundError
	}
	return nil
}

func (repository *projectPermissionRepository) Get(ctx context.Context, dbExecutor database.DatabaseExecutor, permissionID string) (domain.ProjectPermission, error) {
	var dto projectPermissionDTO
	err := dbExecutor.QueryRowContext(
		ctx,
		`
            SELECT
                PUR.id, PUR.user_id, U.name, PUR.role_id, R.name, R.permissions_bitmask
            FROM project_user_role PUR
			INNER JOIN users U ON U.id = PUR.user_id
			INNER JOIN roles R ON R.id = PUR.role_id
            WHERE
				PUR.id = ?
        `,
		permissionID).Scan(&dto.ID, &dto.UserID, &dto.UserName, &dto.RoleID, &dto.RoleName, &dto.RolePermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectPermission{}, domain.NotFoundError
		}
		return domain.ProjectPermission{}, err
	}
	return toDomain(dto), err
}

func (repository *projectPermissionRepository) GetProjectPermissions(ctx context.Context, dbExecutor database.DatabaseExecutor, projectID string) ([]domain.ProjectPermission, error) {
	rows, err := dbExecutor.QueryContext(
		ctx,
		`
            SELECT
                PUR.id, PUR.user_id, U.name, PUR.role_id, R.name, R.permissions_bitmask
            FROM project_user_role PUR
			INNER JOIN users U ON U.id = PUR.user_id
			INNER JOIN roles R ON R.id = PUR.role_id
            WHERE
				PUR.project_id = ?
			ORDER BY
				U.name COLLATE NOCASE
        `,
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]projectPermissionDTO, 0)
	for rows.Next() {
		var dto projectPermissionDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserID, &dto.UserName, &dto.RoleID, &dto.RoleName, &dto.RolePermissionsBitmask,
		); err != nil {
			return nil, err
		}
		dtos = append(dtos, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return toDomainArray(dtos), nil
}
