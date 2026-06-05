package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/aportela/doneo/internal/utils"
)

func createProjectTypes(database database.Database) []string {
	projectTypeNames := []string{
		"Personal", "Business", "Work", "Educational", "Technology",
		"Creative", "Research", "Social", "Marketing", "Sports",
		"Health", "Sustainability", "Government", "Financial", "Construction",
		"Legal", "Logistics", "Administrative", "Strategy",
	}
	var newProjectTypeIds []string
	service := projecttypeservice.NewService(database, projecttyperepository.NewRepository(database))
	for _, projectTypeName := range projectTypeNames {
		projectTypeID := utils.UUID()
		err := service.Add(context.Background(), domain.ProjectType{
			ID:       projectTypeID,
			Name:     projectTypeName,
			HexColor: utils.RandomSoftHexColor(),
		})
		if err != nil {
			fmt.Printf("Error creating project type %s\n", err.Error())
		} else {
			newProjectTypeIds = append(newProjectTypeIds, projectTypeID)
		}
	}
	return newProjectTypeIds
}
