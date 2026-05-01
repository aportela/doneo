package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func createProjectPriorities(database database.Database) []string {
	projectPriorityNames := []string{"Low", "Medium", "High"}
	var newProjectPriorityIds []string
	projectPriorityRepository := projectpriorityrepository.NewProyectPriorityRepository(database)
	projectPriorityService := projectpriorityservice.NewProjectPriorityService(projectPriorityRepository)
	for index, projectPriorityName := range projectPriorityNames {
		projectPriorityID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := projectPriorityService.AddProjectPriority(context.Background(), domain.ProjectPriority{
			ID:       projectPriorityID,
			Name:     projectPriorityName,
			Index:    index,
			HexColor: utils.RandomSoftHexColor(),
		})
		if err != nil {
			fmt.Printf("Error creating project priority %s\n", err.Error())
		}
		newProjectPriorityIds = append(newProjectPriorityIds, projectPriorityID)
	}
	return newProjectPriorityIds
}
