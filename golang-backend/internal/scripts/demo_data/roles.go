package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/roleservice"
)

func createRoles(db database.Database) []string {
	var newRoleIds []string
	authorizationService := authorizationservice.NewService(db, cache.NewPermissionCache(), userrepository.NewRepository(), projectpermissionrepository.NewRepository())
	roleService := roleservice.NewService(db, authorizationService, rolerepository.NewRepository())
	permissionBitMask := domain.Bitmask(0)
	permissionBitMask.AddFlag(domain.PermissionUpdateProject | domain.PermissionDeleteProject | domain.PermissionViewProject | domain.PermissionAddTask | domain.PermissionUpdateTask | domain.PermissionDeleteTask | domain.PermissionViewTask)
	role := domain.Role{
		RoleBase: domain.RoleBase{
			Name: "Administrator",
		},
		PermissionsBitmask: permissionBitMask,
	}
	ctx := middlewares.SetContextUser(context.Background(), middlewares.ContextUser{UserBase: domain.UserBase{}, SkipAuthorization: true})
	role, err := roleService.Add(ctx, role)
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, role.ID)
	}
	permissionBitMask.Clear()
	permissionBitMask.AddFlag(domain.PermissionViewProject | domain.PermissionViewTask)
	role.Name = "Guest"
	role.PermissionsBitmask = permissionBitMask
	role, err = roleService.Add(ctx, role)
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, role.ID)
	}
	return newRoleIds
}
