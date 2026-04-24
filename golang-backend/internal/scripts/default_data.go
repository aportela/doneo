package scripts

import (
	"context"
	"fmt"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/repositories"
	"github.com/aportela/gotask/internal/services"
	"github.com/aportela/gotask/internal/utils"
	"github.com/gofrs/uuid"
)

func CreateDefaultData(db database.Database) {
	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)

	userID := func() string { u, _ := uuid.NewV7(); return u.String() }()
	password := "secret"
	err := userService.AddUser(context.Background(), models.User{
		UserBase: models.UserBase{
			ID:   userID,
			Name: "administrator",
		},
		Email:        "admin@localhost",
		Password:     &password,
		CreatedAt:    utils.CurrentTimestamp(),
		LastUpdateAt: nil,
	})
	if err != nil {
		fmt.Printf("Error creating user %s\n", err.Error())
	}
}
