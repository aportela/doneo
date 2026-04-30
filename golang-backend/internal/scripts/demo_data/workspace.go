package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/repositories/workspacerepository"
	"github.com/aportela/doneo/internal/services/workspaceservice"
	"github.com/aportela/doneo/internal/utils"
	"github.com/gofrs/uuid"
)

func createDefaultWorkspace(database database.Database, userId string) string {
	workspaceRepository := workspacerepository.NewWorkspaceRepository(database)
	workspaceService := workspaceservice.NewWorkspaceService(workspaceRepository)
	workspaceID := func() string { u, _ := uuid.NewV7(); return u.String() }()
	description := "Default workspace"
	err := workspaceService.AddWorkspace(context.Background(), domain.Workspace{
		ID:          workspaceID,
		Name:        "Default",
		Description: &description,
		CreatedBy:   domain.UserBase{ID: userId},
		CreatedAt:   utils.CurrentMSTimestamp(),
	})
	if err != nil {
		fmt.Printf("Error creating workspace %s\n", err.Error())
	}
	return workspaceID
}
