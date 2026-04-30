package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/gofrs/uuid"
)

func createProjectStatuses(database database.Database) []string {
	projectStatusNames := []string{
		"Pending", "Started", "Stopped", "Finished", "Aborted",
	}
	var newProjectStatusIds []string
	projectStatusRepository := projectstatusrepository.NewProjectStatusRepository(database)
	projectStatusService := projectstatusservice.NewProjectStatusService(projectStatusRepository)
	for index, projectStatusName := range projectStatusNames {
		projectStatusID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := projectStatusService.AddProjectStatus(context.Background(), domain.ProjectStatus{
			ID:    projectStatusID,
			Name:  projectStatusName,
			Index: index,
		})
		if err != nil {
			fmt.Printf("Error creating project status %s\n", err.Error())
		}
		newProjectStatusIds = append(newProjectStatusIds, projectStatusID)
	}
	return newProjectStatusIds
}
