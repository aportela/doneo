package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/services/roleservice"
)

func createRoles(database database.Database) []string {
	var newRoleIds []string
	roleService := roleservice.NewService(database, rolerepository.NewRepository(database))
	permissionBitMask := domain.Bitmask(0)
	permissionBitMask.AddFlag(domain.PermissionUpdateProject | domain.PermissionDeleteProject | domain.PermissionViewProject | domain.PermissionAddTask | domain.PermissionUpdateTask | domain.PermissionDeleteTask | domain.PermissionViewTask)
	role := domain.Role{
		RoleBase: domain.RoleBase{
			Name: "Administrator",
		},
		PermissionsBitmask: permissionBitMask,
	}
	role, err := roleService.Add(context.Background(), role)
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, role.ID)
	}
	permissionBitMask.Clear()
	permissionBitMask.AddFlag(domain.PermissionViewProject | domain.PermissionViewTask)
	role.Name = "Guest"
	role.PermissionsBitmask = permissionBitMask
	role, err = roleService.Add(context.Background(), role)
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, role.ID)
	}
	return newRoleIds
}
