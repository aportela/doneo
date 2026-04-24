package seed

import (
	"github.com/aportela/gotask/internal/database"
)

func CreateDemoData(db database.Database) {
	// TODO
	/*
		projectRepository := repositories.NewProjectRepository(db)
		projectService := services.NewProjectService(projectRepository)
		err := projectService.AddProject(context.Background(), models.Project{
			ID:             func() string { u, _ := uuid.NewV7(); return u.String() }(),
			Key:            "AA",
			Summary:        "Summary",
			Description:    nil,
			CreatedBy:      models.UserBase{ID: "019dba5d-83a4-7f97-bdf1-97a5fb3d5869"},
			CreatedAt:      utils.CurrentTimestamp(),
			LastModifiedAt: nil,
			StartedAt:      nil,
			FinishedAt:     nil,
			DueAt:          nil,
			Type:           models.ProjectType{ID: "019dba85-0669-7fd4-86ed-dbe36df285af"},
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	*/

}
