package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/userservice"
)

func CreateDefaultAdminUser(db database.Database) {
	permissionsBitmask := domain.Bitmask(0)
	permissionsBitmask.AddFlag(domain.UserPermissionAdmin)
	service := userservice.NewService(db, userrepository.NewRepository(db))
	_, err := service.Add(context.Background(), domain.User{
		UserBase: domain.UserBase{
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
