package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/utils"
)

func createProjectPriorities(database database.Database) []string {
	projectPriorityNames := []string{"Low", "Medium", "High"}
	var newProjectPriorityIds []string
	service := projectpriorityservice.NewService(database, projectpriorityrepository.NewRepository(database))
	for index, projectPriorityName := range projectPriorityNames {
		projectPriority := domain.ProjectPriority{
			Name:     projectPriorityName,
			HexColor: utils.RandomSoftHexColor(),
			Index:    uint(index),
		}
		projectPriority, err := service.Add(context.Background(), projectPriority)
		if err != nil {
			fmt.Printf("Error creating project priority %s\n", err.Error())
		}
		newProjectPriorityIds = append(newProjectPriorityIds, projectPriority.ID)
	}
	return newProjectPriorityIds
}
