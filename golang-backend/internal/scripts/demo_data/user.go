package demodatascripts

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/userservice"
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
	email := fmt.Sprintf("%s%s%d@localhost.local", firstInitial, secondAndThirdName, randomNumber)
	return email
}

func getRandomUser() domain.User {
	name := getRandomUserName()
	permissionsBitmask := domain.Bitmask(0)
	//permissionsBitmask.AddPermission(domain.UserPermissionAdmin)
	return domain.User{
		UserBase: domain.UserBase{
			Name: name,
		},
		Email:              generateRandomEmail(name),
		CreatedAt:          time.Now(),
		UpdatedAt:          nil,
		DeletedAt:          nil,
		PermissionsBitmask: permissionsBitmask,
	}
}

func createUsers(db database.Database, count int) []string {
	var newUserIds []string
	service := userservice.NewService(db, userrepository.NewRepository(db))
	for i := 1; i <= count; i++ {
		newUser := getRandomUser()
		newUser, err := service.Add(context.Background(), newUser, "secret")
		if err != nil {
			fmt.Printf("Error creating user %s\n", err.Error())
		} else {
			newUserIds = append(newUserIds, newUser.ID)
		}
	}
	return newUserIds
}
