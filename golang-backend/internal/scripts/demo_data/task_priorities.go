package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func createTaskPriorities(database database.Database) []string {
	taskPriorityNames := []string{"Low", "Medium", "High"}
	var newTaskPriorityIds []string
	taskPriorityService := taskpriorityservice.NewService(database, taskpriorityrepository.NewRepository(database))
	for _, taskPriorityName := range taskPriorityNames {
		taskPriorityID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := taskPriorityService.Add(context.Background(), domain.TaskPriority{
			ID:       taskPriorityID,
			Name:     taskPriorityName,
			HexColor: utils.RandomSoftHexColor(),
		})
		if err != nil {
			fmt.Printf("Error creating task priority %s\n", err.Error())
		}
		newTaskPriorityIds = append(newTaskPriorityIds, taskPriorityID)
	}
	return newTaskPriorityIds
}
