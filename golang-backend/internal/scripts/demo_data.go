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

func CreateDemoData(db database.Database) {
	projectRepository := repositories.NewProjectRepository(db)
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
			CreatedAt:      utils.CurrentTimestamp(),
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
