package scripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/models"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
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
		Email:           "admin@localhost",
		Password:        &password,
		CreatedAt:       utils.CurrentMSTimestamp(),
		LastUpdateAt:    nil,
		IsAdministrator: true,
	})
	if err != nil {
		fmt.Printf("Error creating user %s\n", err.Error())
	}
}
