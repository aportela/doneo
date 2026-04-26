package demodatascripts

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/models"
	"github.com/aportela/doneo/internal/repositories"
	"github.com/aportela/doneo/internal/services"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func getRandomUserName() string {
	names := []string{
		"James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Charles", "Thomas",
		"Mary", "Jennifer", "Linda", "Patricia", "Elizabeth", "Susan", "Jessica", "Sarah", "Karen", "Nancy",
	}

	surnames := []string{
		"Smith", "Johnson", "Williams", "Jones", "Brown", "Davis", "Miller", "Wilson", "Moore", "Taylor",
		"Anderson", "Thomas", "Jackson", "White", "Harris", "Martin", "Thompson", "Garcia", "Martinez", "Roberts",
	}

	name := names[rand.Intn(len(names))]
	surname1 := surnames[rand.Intn(len(surnames))]
	surname2 := surnames[rand.Intn(len(surnames))]

	return name + " " + surname1 + " " + surname2
}

func generateRandomEmail(fullName string) string {
	parts := strings.Fields(fullName)
	firstInitial := strings.ToLower(string(parts[0][0]))
	secondAndThirdName := strings.ToLower(parts[1]) + "." + strings.ToLower(parts[2])
	randomNumber := rand.Intn(100)
	email := fmt.Sprintf("%s%s%d@localhost", firstInitial, secondAndThirdName, randomNumber)
	return email
}

func getRandomUser() models.User {
	userID := func() string { u, _ := uuid.NewV7(); return u.String() }()
	password := userID
	name := getRandomUserName()
	return models.User{
		UserBase: models.UserBase{
			ID:   userID,
			Name: name,
		},
		Email:       generateRandomEmail(name),
		Password:    &password,
		CreatedAt:   utils.GetRandomMSTimestamp(time.Now().AddDate(-5, 0, 0), time.Now()),
		UpdatedAt:   nil,
		IsSuperUser: false,
	}
}

func createUsers(database database.Database, count int) []string {
	var newUserIds []string
	userRepository := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepository)
	for i := 1; i <= count; i++ {
		newUser := getRandomUser()
		err := userService.AddUser(context.Background(), newUser)
		if err != nil {
			fmt.Printf("Error creating user %s\n", err.Error())
		}
		newUserIds = append(newUserIds, newUser.ID)
	}
	return newUserIds
}
