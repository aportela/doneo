package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/utils"
)

func createProjectStatuses(db database.Database) []string {
	projectStatusNames := []string{
		"Pending", "Started", "Stopped", "Finished", "Aborted",
	}
	var newProjectStatusIds []string
	authorizationService := authorizationservice.NewService(db, cache.NewPermissionCache(), userrepository.NewRepository(), projectpermissionrepository.NewRepository())
	service := projectstatusservice.NewService(db, authorizationService, projectstatusrepository.NewRepository())
	for index, projectStatusName := range projectStatusNames {
		var flags domain.Bitmask
		switch projectStatusName {
		case "Pending":
			flags.AddFlag(domain.ProjectStatusFlagDefaultOnCreate)
		case "Started":
			flags.AddFlag(domain.ProjectStatusFlagFillEmptyStartDate)
		case "Finished":
			flags.AddFlag(domain.ProjectStatusFlagFillEmptyFinishDate)
			flags.AddFlag(domain.ProjectStatusFlagUnsetFinishDateOnLeave)
		}
		projectStatus := domain.ProjectStatus{
			Name:     projectStatusName,
			HexColor: utils.RandomSoftHexColor(),
			Index:    uint(index),
			Flags:    flags,
		}
		projectStatus, err := service.Add(context.Background(), projectStatus)
		if err != nil {
			fmt.Printf("Error creating project status %s\n", err.Error())
		}
		newProjectStatusIds = append(newProjectStatusIds, projectStatus.ID)
	}
	return newProjectStatusIds
}
