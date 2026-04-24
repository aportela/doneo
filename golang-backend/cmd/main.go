package main

import (
	"log"
	"net/http"

	"github.com/aportela/gotask/internal/cli"
	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/router"
)

func main() {
	log.Println("starting GOTask v0.1alpha...")

	db, err := database.Open(true)
	if err != nil {
		log.Fatal(err)
	} else {

		//db.Close()

		params, err := cli.HandleFlags()
		if err != nil {
			log.Fatal(err)
		}

		if params.InsertBulkData {
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

		r := router.NewRouter(db)

		log.Println("Listening over http://localhost:3000/")

		http.ListenAndServe(":3000", r)
	}
}
