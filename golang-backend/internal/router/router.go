package router

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/gotask/internal/fileserver"
	"github.com/aportela/gotask/internal/handlers"
	"github.com/aportela/gotask/internal/repositories"
	"github.com/aportela/gotask/internal/services"
)

func NewRouter(db *sql.DB) http.Handler {
	baseRouter := chi.NewRouter()

	baseRouter.Use(middleware.Logger)

	apiRouter := chi.NewRouter()
	projectRepository := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepository)
	projectHandler := handlers.NewProjectHandler(projectService)
	apiRouter.Get("/project/{id}", projectHandler.GetProject)
	apiRouter.Get("/projects", projectHandler.SearchProjects)
	baseRouter.Mount("/api", apiRouter)

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	publicWebPath := http.Dir(filepath.Join(workDir, "public"))
	fileserver.FileServer(baseRouter, "/", publicWebPath)

	return baseRouter
}
