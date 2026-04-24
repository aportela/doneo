package router

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/fileserver"
	"github.com/aportela/gotask/internal/handlers"
	"github.com/aportela/gotask/internal/repositories"
	"github.com/aportela/gotask/internal/services"
)

func NewRouter(db database.Database) http.Handler {
	baseRouter := chi.NewRouter()

	baseRouter.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/users", func(r chi.Router) {
		userRepository := repositories.NewUserRepository(db)
		userService := services.NewUserService(userRepository)
		userHandler := handlers.NewUserHandler(userService)
		r.Post("/", userHandler.AddUser)
		r.Get("/", userHandler.SearchUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/", userHandler.DeleteUser)
		})
	})

	apiRouter.Route("/projects", func(r chi.Router) {
		projectRepository := repositories.NewProjectRepository(db)
		projectService := services.NewProjectService(projectRepository)
		projectHandler := handlers.NewProjectHandler(projectService)
		r.Post("/", projectHandler.AddProject)
		r.Get("/", projectHandler.SearchProjects)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", projectHandler.GetProject)
			r.Put("/", projectHandler.UpdateProject)
			r.Delete("/", projectHandler.DeleteProject)
		})
	})

	apiRouter.Route("/project_types", func(r chi.Router) {
		projectTypeRepository := repositories.NewProjectTypeRepository(db)
		projectTypeService := services.NewProjectTypeService(projectTypeRepository)
		projectTypeHandler := handlers.NewProjectTypeHandler(projectTypeService)
		r.Post("/", projectTypeHandler.AddProjectType)
		r.Get("/", projectTypeHandler.SearchProjectTypes)
		r.Route("/{id}", func(r chi.Router) {
			r.Put("/", projectTypeHandler.UpdateProjectType)
			r.Delete("/", projectTypeHandler.DeleteProjectType)
		})
	})

	baseRouter.Mount("/api", apiRouter)

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	publicWebPath := http.Dir(filepath.Join(workDir, "public"))
	fileserver.FileServer(baseRouter, "/", publicWebPath)

	return baseRouter
}
