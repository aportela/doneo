package router

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/doneo/internal/configuration"
	"github.com/aportela/doneo/internal/database"
	authHandlers "github.com/aportela/doneo/internal/handlers/auth"
	userHandlers "github.com/aportela/doneo/internal/handlers/user"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/ui"
)

func NewRouter(db database.Database, config configuration.Configuration) http.Handler {
	baseRouter := chi.NewRouter()

	baseRouter.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/auth", func(r chi.Router) {
		userHandler := authHandlers.NewAuthHandler(db, config.Auth.SecretKey, config.Auth.AccessTokenExpirationDays, config.Auth.RefreshTokenExpirationDays)
		r.Post("/signup", userHandler.SignUp)
		r.Post("/signin", userHandler.SignIn)
		r.Post("/signout", userHandler.SignOut)
		r.Post("/renew_access_token", userHandler.RenewAccessToken)

	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Use(middlewares.RequireAuthentication)
		baseRouter.Use(middlewares.RequireSuperUser)
		userHandler := userHandlers.NewUserHandler(db)
		r.Post("/", userHandler.AddUser)
		r.Get("/", userHandler.SearchUsers)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", userHandler.GetUser)
			r.Put("/", userHandler.UpdateUser)
			r.Delete("/", userHandler.DeleteUser)
		})
	})

	/*
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

	*/

	baseRouter.Mount("/api", apiRouter)

	subFS, err := fs.Sub(ui.Dist, "dist")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(subFS))
	baseRouter.Handle("/*", fileServer)

	baseRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		data, _ := fs.ReadFile(subFS, "index.html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(data)
	})

	return baseRouter
}
