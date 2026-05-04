package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/userservice"
	"github.com/gofrs/uuid"
)

func CreateDefaultData(db database.Database) {
	userRepository := userrepository.NewUserRepository(db)
	userService := userservice.NewUserService(userRepository)

	userID := func() string { u, _ := uuid.NewV7(); return u.String() }()
	err := userService.AddUser(context.Background(), domain.User{
		UserBase: domain.UserBase{
			ID:   userID,
			Name: "administrator",
		},
		Email:       "admin@localhost.localnet",
		Password:    "secret",
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		DeletedAt:   nil,
		IsSuperUser: true,
	})
	if err != nil {
		fmt.Printf("Error creating user %s\n", err.Error())
	}

	/*
		workspaceRepository := workspacerepository.NewWorkspaceRepository(db)
		workspaceService := workspaceservice.NewWorkspaceService(workspaceRepository)

		workspaceID := func() string { u, _ := uuid.NewV7(); return u.String() }()

		err = workspaceService.AddWorkspace(context.Background(), domain.Workspace{
			ID:          workspaceID,
			Name:        "default",
			Description: nil,
			HexColor:    utils.RandomSoftHexColor(),
			CreatedBy: domain.UserBase{
				ID: userID,
			},
			CreatedAt: utils.CurrentMSTimestamp(),
			UpdatedAt: nil,
		})
		if err != nil {
			fmt.Printf("Error creating workspace %s\n", err.Error())
		}
	*/
}
