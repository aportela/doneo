package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/userservice"
	"github.com/aportela/doneo/internal/utils"
)

func CreateDefaultAdminUser(database database.Database) {
	permissionsBitmask := domain.PermissionsBitmask(0)
	permissionsBitmask.AddPermission(domain.UserPermissionAdmin)
	service := userservice.NewService(database, userrepository.NewRepository(database))
	err := service.Add(context.Background(), domain.User{
		UserBase: domain.UserBase{
			ID:   utils.UUID(),
			Name: "administrator",
		},
		Email:              "admin@localhost.localnet",
		CreatedAt:          time.Now(),
		UpdatedAt:          nil,
		DeletedAt:          nil,
		PermissionsBitmask: permissionsBitmask,
	}, "secret")
	if err != nil {
		fmt.Printf("Error creating user %s\n", err.Error())
	}
}
