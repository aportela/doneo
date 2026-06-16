package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/utils"
)

func createTaskPriorities(db database.Database) []string {
	taskPriorityNames := []string{"Low", "Medium", "High"}
	var newTaskPriorityIds []string
	authorizationService := authorizationservice.NewService(db, cache.NewPermissionCache(), userrepository.NewRepository(), projectpermissionrepository.NewRepository())
	taskPriorityService := taskpriorityservice.NewService(db, authorizationService, taskpriorityrepository.NewRepository())
	for index, taskPriorityName := range taskPriorityNames {
		taskPriority := domain.TaskPriority{
			Name:     taskPriorityName,
			HexColor: utils.RandomSoftHexColor(),
			Index:    uint(index),
		}
		taskPriority, err := taskPriorityService.Add(context.Background(), taskPriority)
		if err != nil {
			fmt.Printf("Error creating task priority %s\n", err.Error())
		}
		newTaskPriorityIds = append(newTaskPriorityIds, taskPriority.ID)
	}
	return newTaskPriorityIds
}
