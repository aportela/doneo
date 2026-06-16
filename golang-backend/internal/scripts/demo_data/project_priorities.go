package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/utils"
)

func createProjectPriorities(db database.Database) []string {
	projectPriorityNames := []string{"Low", "Medium", "High"}
	var newProjectPriorityIds []string
	authorizationService := authorizationservice.NewService(db, cache.NewPermissionCache(), userrepository.NewRepository(), projectpermissionrepository.NewRepository())
	service := projectpriorityservice.NewService(db, authorizationService, projectpriorityrepository.NewRepository())
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
