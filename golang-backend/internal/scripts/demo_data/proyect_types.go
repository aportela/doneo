package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/models"
	"github.com/aportela/gotask/internal/repositories"
	"github.com/aportela/gotask/internal/services"
	"github.com/gofrs/uuid"
)

func createProjectTypes(database database.Database) []string {

	projectTypeNames := []string{
		"Personal", "Business", "Work", "Educational", "Technology",
		"Creative", "Research", "Social", "Marketing", "Sports",
		"Health", "Sustainability", "Government", "Financial", "Construction",
		"Legal", "Logistics", "Administrative", "Strategy",
	}

	var newProjectTypeIds []string
	projectTypeRepository := repositories.NewProjectTypeRepository(database)
	projectTypeService := services.NewProjectTypeService(projectTypeRepository)
	for _, projectTypeName := range projectTypeNames {
		projectTypeID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := projectTypeService.AddProjectType(context.Background(), models.ProjectType{

			ID:   projectTypeID,
			Name: projectTypeName,
		})
		if err != nil {
			fmt.Printf("Error creating project type %s\n", err.Error())
		}
		newProjectTypeIds = append(newProjectTypeIds, projectTypeID)
	}
	return newProjectTypeIds
}
