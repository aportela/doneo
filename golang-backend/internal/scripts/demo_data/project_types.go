package demodatascripts

import (
	"context"
	"fmt"

	"github.com/aportela/doneo/internal/cache"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authorizationservice"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/aportela/doneo/internal/utils"
)

func createProjectTypes(db database.Database) []string {
	projectTypeNames := []string{
		"Personal", "Business", "Work", "Educational", "Technology",
		"Creative", "Research", "Social", "Marketing", "Sports",
		"Health", "Sustainability", "Government", "Financial", "Construction",
		"Legal", "Logistics", "Administrative", "Strategy",
	}
	var newProjectTypeIds []string
	authorizationService := authorizationservice.NewService(db, cache.NewPermissionCache(), userrepository.NewRepository(), projectpermissionrepository.NewRepository())
	service := projecttypeservice.NewService(db, authorizationService, projecttyperepository.NewRepository())
	ctx := middlewares.SetContextUser(context.Background(), middlewares.ContextUser{UserBase: domain.UserBase{}, SkipAuthorization: true})
	for _, projectTypeName := range projectTypeNames {
		projectType := domain.ProjectType{
			Name:     projectTypeName,
			HexColor: utils.RandomSoftHexColor(),
		}
		projectType, err := service.Add(ctx, projectType)
		if err != nil {
			fmt.Printf("Error creating project type %s\n", err.Error())
		} else {
			newProjectTypeIds = append(newProjectTypeIds, projectType.ID)
		}
	}
	return newProjectTypeIds
}
