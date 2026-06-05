package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func createProjectStatuses(database database.Database) []string {
	projectStatusNames := []string{
		"Pending", "Started", "Stopped", "Finished", "Aborted",
	}
	var newProjectStatusIds []string
	service := projectstatusservice.NewService(database, projectstatusrepository.NewRepository(database))
	for _, projectStatusName := range projectStatusNames {
		projectStatusID := func() string { u, _ := uuid.NewV7(); return u.String() }()
		err := service.Add(context.Background(), domain.ProjectStatus{
			ID:       projectStatusID,
			Name:     projectStatusName,
			HexColor: utils.RandomSoftHexColor(),
		})
		if err != nil {
			fmt.Printf("Error creating project status %s\n", err.Error())
		}
		newProjectStatusIds = append(newProjectStatusIds, projectStatusID)
	}
	return newProjectStatusIds
}
