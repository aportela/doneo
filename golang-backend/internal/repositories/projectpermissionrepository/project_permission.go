package projectpermissionrepository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectPermissionRepository interface {
	Add(ctx context.Context, projectId string, permission domain.ProjectPermission) error
	Delete(ctx context.Context, projectId string, permissionId string) error
	Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error)
	Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error)
}

type projectPermissionRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectPermissionRepository {
	return &projectPermissionRepository{database: database}
}

func (repository *projectPermissionRepository) Add(ctx context.Context, projectId string, permission domain.ProjectPermission) error {
	dto := toDTO(permission)
	_, err := repository.database.ExecContext(
		ctx,
		`
            INSERT INTO project_user_role (id, project_id, user_id, role_id)
			VALUES (?, ?, ?, ?)
        `,
		dto.ID,
		projectId,
		dto.UserId,
		dto.RoleId,
	)
	if err != nil {
		// TODO: remove ?
		fmt.Println(err.Error())
		var sqlErr *sqlite.Error
		if !errors.As(err, &sqlErr) {
			return err
		}
		switch sqlErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return &domain.ValidationError{Field: "projectId,userId,roleId"}
		case sqlite3.SQLITE_CONSTRAINT_CHECK:
			if strings.Contains(sqlErr.Error(), "length(project_id)") {
				return &domain.ValidationError{Field: "projectId"}
			} else if strings.Contains(sqlErr.Error(), "length(user_id)") {
				return &domain.ValidationError{Field: "userId"}
			} else if strings.Contains(sqlErr.Error(), "length(role_id)") {
				return &domain.ValidationError{Field: "roleId"}
			}
			return err
		default:
			return err
		}
	}
	return nil
}

func (repository *projectPermissionRepository) Delete(ctx context.Context, projectId string, permissionId string) error {
	_, err := repository.database.ExecContext(
		ctx,
		`
            DELETE FROM project_user_role
			WHERE
				id = ?
        `,
		permissionId,
	)
	if err != nil {
		// TODO: remove ?
		// TODO: check sql error
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (repository *projectPermissionRepository) Get(ctx context.Context, permissionId string) (domain.ProjectPermission, error) {
	var dto projectPermissionDTO
	err := repository.database.QueryRowContext(
		ctx,
		`
            SELECT
                PUR.id, PUR.user_id, U.name, PUR.role_id, R.name, R.permissions_bitmask
            FROM project_user_role PUR
			INNER JOIN users U ON U.id = PUR.user_id
			INNER JOIN roles R ON R.id = PUR.role_id
            WHERE PUR.id = ?
        `,
		permissionId).Scan(&dto.ID, &dto.UserId, &dto.UserName, &dto.RoleId, &dto.RoleName, &dto.RolePermissionsBitmask)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ProjectPermission{}, domain.NotFoundError
		}
		return domain.ProjectPermission{}, err
	}
	return toDomain(dto), err
}

func (repository *projectPermissionRepository) Search(ctx context.Context, projectId string) ([]domain.ProjectPermission, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
                PUR.id, PUR.user_id, U.name, PUR.role_id, R.name, R.permissions_bitmask
            FROM project_user_role PUR
			INNER JOIN users U ON U.id = PUR.user_id
			INNER JOIN roles R ON R.id = PUR.role_id
            WHERE PUR.project_id = ?
			ORDER BY U.name COLLATE NOCASE
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	dtos := make([]projectPermissionDTO, 0)
	for rows.Next() {
		var dto projectPermissionDTO
		if err := rows.Scan(
			&dto.ID, &dto.UserId, &dto.UserName, &dto.RoleId, &dto.RoleName, &dto.RolePermissionsBitmask,
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
