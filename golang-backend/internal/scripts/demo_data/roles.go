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
	permissionBitMask.AddFlag(domain.PermissionCreate | domain.PermissionUpdate | domain.PermissionDelete | domain.PermissionView | domain.PermissionList | domain.PermissionExecute)
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
	permissionBitMask.AddFlag(domain.PermissionView | domain.PermissionList)
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
