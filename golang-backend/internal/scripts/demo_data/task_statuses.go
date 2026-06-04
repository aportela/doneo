package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func createTaskStatuses(database database.Database) []string {
	taskStatusNames := []string{
		"Pending", "Started", "Stopped", "Finished", "Aborted",
	}
	var newTaskStatusIds []string
	projectStatusService := taskstatusservice.NewService(database, taskstatusrepository.NewRepository(database))
	for _, taskStatusName := range taskStatusNames {
		taskStatusID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := projectStatusService.Add(context.Background(), domain.TaskStatus{
			ID:       taskStatusID,
			Name:     taskStatusName,
			HexColor: utils.RandomSoftHexColor(),
		})
		if err != nil {
			fmt.Printf("Error creating tas status %s\n", err.Error())
		}
		newTaskStatusIds = append(newTaskStatusIds, taskStatusID)
	}
	return newTaskStatusIds
}
