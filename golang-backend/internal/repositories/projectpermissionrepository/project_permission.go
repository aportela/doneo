package projectpermissionrepository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type ProjectPermissionRepository interface {
	Add(ctx context.Context, permissionId string, projectId string, userId string, roleId string) error
	Delete(ctx context.Context, projectId string, permissionId string) error
	Search(ctx context.Context, projectId string) ([]projectPermissionDTO, error)
}

type projectPermissionRepository struct {
	database database.Database
}

func NewRepository(database database.Database) ProjectPermissionRepository {
	return &projectPermissionRepository{database: database}
}

func (repository *projectPermissionRepository) Add(ctx context.Context, permissionId string, projectId string, userId string, roleId string) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(
		ctx,
		`
            INSERT INTO project_user_role (id, project_id, user_id, role_id)
			VALUES (?, ?, ?, ?)
        `,
		permissionId,
		projectId,
		userId,
		roleId,
	)
	if err != nil {
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
		}
	}
	userId, _ = middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectPermissionAdded,
		userId,
		time.Now().UnixMilli(),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *projectPermissionRepository) Delete(ctx context.Context, projectId string, permissionId string) error {
	tx, err := repository.database.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	// TODO: add project id into WHERE ?
	_, err = tx.ExecContext(
		ctx,
		`
            DELETE FROM project_user_role
			WHERE
				id = ?
        `,
		permissionId,
	)
	userId, _ := middlewares.GetUserIDFromContext(ctx)
	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO project_history_operations
				(project_id, operation_type, user_id, created_at)
			VALUES
				(?, ?, ?, ?)
		`,
		projectId,
		domain.EventProjectPermissionDeleted,
		userId,
		time.Now().UnixMilli(),
	)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return tx.Commit()
}

func (repository *projectPermissionRepository) Search(ctx context.Context, projectId string) ([]projectPermissionDTO, error) {
	rows, err := repository.database.QueryContext(
		ctx,
		`
            SELECT
                PUR.id, PUR.user_id, U.name, PUR.role_id, R.name, R.permissions_bitmask
            FROM project_user_role PUR
			INNER JOIN users U ON U.id = PUR.user_id
			INNER JOIN roles R ON R.id = PUR.role_id
            WHERE PUR.project_id = ?
			ORDER BY U.name
        `,
		projectId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	projectPermissions := make([]projectPermissionDTO, 0)
	for rows.Next() {
		var projectPermission projectPermissionDTO
		if err := rows.Scan(
			&projectPermission.ID, &projectPermission.UserId, &projectPermission.UserName, &projectPermission.RoleId, &projectPermission.RoleName, &projectPermission.PermissionsBitmask,
		); err != nil {
			return nil, err
		}
		projectPermissions = append(projectPermissions, projectPermission)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return projectPermissions, nil
}
