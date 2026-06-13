package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/utils"
)

func createTaskStatuses(database database.Database) []string {
	taskStatusNames := []string{
		"Pending", "Started", "Stopped", "Finished", "Aborted",
	}
	var newTaskStatusIds []string
	projectStatusService := taskstatusservice.NewService(database, taskstatusrepository.NewRepository(database))
	for index, taskStatusName := range taskStatusNames {
		var flags domain.Bitmask
		switch taskStatusName {
		case "Pending":
			flags.AddFlag(domain.TaskStatusFlagDefaultOnCreate)
		case "Started":
			flags.AddFlag(domain.TaskStatusFlagFillEmptyStartDate)
		case "Finished":
			flags.AddFlag(domain.TaskStatusFlagFillEmptyFinishDate)
			flags.AddFlag(domain.TaskStatusFlagUnsetFinishDateOnLeave)
		}
		taskStatus := domain.TaskStatus{
			Name:     taskStatusName,
			HexColor: utils.RandomSoftHexColor(),
			Index:    uint(index),
		}
		taskStatus, err := projectStatusService.Add(context.Background(), taskStatus)
		if err != nil {
			fmt.Printf("Error creating tas status %s\n", err.Error())
		}
		newTaskStatusIds = append(newTaskStatusIds, taskStatus.ID)
	}
	return newTaskStatusIds
}
