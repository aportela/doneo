package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/aportela/doneo/internal/utils"
)

func createRoles(database database.Database) []string {
	var newRoleIds []string
	roleRepository := rolerepository.NewRoleRepository(database)
	roleService := roleservice.NewRoleService(roleRepository)
	roleID := utils.UUID()
	var permissionBitMask domain.PermissionsBitmask
	permissionBitMask = 0
	permissionBitMask.AddPermission(domain.PermissionCreate | domain.PermissionUpdate | domain.PermissionDelete | domain.PermissionView | domain.PermissionList | domain.PermissionExecute)
	err := roleService.Add(context.Background(), domain.Role{
		ID:                 roleID,
		Name:               "Administrator",
		PermissionsBitmask: permissionBitMask,
	})
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, roleID)
	}
	roleID = utils.UUID()
	permissionBitMask = 0
	permissionBitMask.AddPermission(domain.PermissionView | domain.PermissionList)
	err = roleService.Add(context.Background(), domain.Role{
		ID:                 roleID,
		Name:               "Guest",
		PermissionsBitmask: permissionBitMask,
	})
	if err != nil {
		fmt.Printf("Error creating role %s\n", err.Error())
	} else {
		newRoleIds = append(newRoleIds, roleID)
	}
	return newRoleIds
}
