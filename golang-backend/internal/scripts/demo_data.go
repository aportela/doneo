package scripts

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/repositories"
	"github.com/aportela/gotask/internal/services"
	"github.com/aportela/gotask/internal/utils"
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

func createUsers(database database.Database, count int) [] string {
	newUserIds []string
	userRepository := repositories.NewUserRepository(database)
	userService := services.NewUserService(userRepository)
	for i := 1; i <= count; i++ {
		userID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		password := userID
		err := userService.AddUser(context.Background(), models.User{
			UserBase: models.UserBase{
				ID:   userID,
				Name: getRandomUserName(),
			},
			Email:           userID + "@localhost",
			Password:        &password,
			CreatedAt:       utils.CurrentMSTimestamp(),
			LastUpdateAt:    nil,
			IsAdministrator: true,
		})
		if err != nil {
			fmt.Printf("Error creating user %s\n", err.Error())
		}
		append(newUserIds, userID)
	}
	return newUserIds
}

func CreateDemoData(database database.Database) {

	userIds := createUsers(database, 32)
	projectRepository := repositories.NewProjectRepository(database)
	projectService := services.NewProjectService(projectRepository)

	projectNames := []string{
		"Customer Portal", "Website Redesign", "Mobile App Launch", "CRM Integration",
		"Marketing Campaign Q2", "Data Analytics Dashboard", "Internal Tools Upgrade",
		"Cloud Migration", "Security Audit", "Product Feature Expansion",
		"Supply Chain Optimization", "Inventory Management System", "Employee Onboarding App",
		"AI Chatbot Integration", "Customer Feedback Portal", "Sales Dashboard Revamp",
		"Brand Awareness Campaign", "Performance Review Automation", "Project Management Tool",
		"Email Marketing Automation", "SEO Optimization", "Payment Gateway Integration",
		"Data Warehouse Upgrade", "Bug Tracking System", "Client Reporting Tool",
		"Business Intelligence Dashboard", "Remote Collaboration Suite", "IT Infrastructure Upgrade",
		"Knowledge Base Platform", "Employee Training Portal", "Social Media Analytics",
		"Website Accessibility Update",
	}

	for i := 1; i <= 128; i++ {
		projectID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		key := fmt.Sprintf("PROJ%02d", i)
		randomName := projectNames[rand.Intn(len(projectNames))]
		summary := fmt.Sprintf("%s #%d", randomName, i)
		var description *string
		if rand.Intn(2) == 0 {
			description = nil
		} else {
			descText := fmt.Sprintf("Description Project %d", i)
			description = &descText
		}

		err := projectService.AddProject(context.Background(), models.Project{
			ID:             projectID,
			Key:            key,
			Summary:        summary,
			Description:    description,
			CreatedBy:      models.UserBase{ID: "019dba5d-83a4-7f97-bdf1-97a5fb3d5869"},
			CreatedAt:      utils.CurrentMSTimestamp(),
			LastModifiedAt: nil,
			StartedAt:      nil,
			FinishedAt:     nil,
			DueAt:          nil,
			Type:           models.ProjectType{ID: "019dba85-0669-7fd4-86ed-dbe36df285af"},
		})
		if err != nil {
			fmt.Printf("Error creating demo project %d: %s\n", i, err.Error())
		}
	}
}
