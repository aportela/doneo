package router

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/doneo/internal/config"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers/authhandler"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/middlewares"

	"github.com/aportela/doneo/internal/ui"
)

func NewRouter(db database.Database, cfg config.Configuration) http.Handler {
	baseRouter := chi.NewRouter()

	baseRouter.Use(middleware.Logger)

	apiRouter := chi.NewRouter()

	apiRouter.Route("/auth", func(r chi.Router) {
		userHandler := authhandler.NewAuthHandler(db, cfg.Auth.SecretKey, cfg.Auth.AccessTokenExpirationHours, cfg.Auth.RefreshTokenExpirationDays)
		r.Post("/signin", userHandler.SignIn)
		r.Post("/signout", userHandler.SignOut)
		r.Post("/renew-access-token", userHandler.RenewAccessToken)
	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		userHandler := userhandler.NewUserHandler(db)
		r.Post("/", userHandler.Add)
		r.Post("/search", userHandler.Search)
		r.Get("/{id}", userHandler.Get)
		r.Put("/{id}", userHandler.Update)
		r.Patch("/{id}", userHandler.Patch)
		r.Delete("/{id}", userHandler.Delete)
	})

	apiRouter.Route("/roles", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		roleHandler := rolehandler.NewRoleHandler(db)
		r.Post("/", roleHandler.Add)
		r.Post("/search", roleHandler.Search)
		r.Get("/{id}", roleHandler.Get)
		r.Put("/{id}", roleHandler.Update)
		r.Delete("/{id}", roleHandler.Delete)
	})

	apiRouter.Route("/project-types", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectTypeHandler := projecttypehandler.NewProjectTypeHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Post("/search", projectTypeHandler.Search)
		r.Get("/{id}", projectTypeHandler.Get)
		r.Put("/{id}", projectTypeHandler.Update)
		r.Delete("/{id}", projectTypeHandler.Delete)
	})

	apiRouter.Route("/project-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectStatusHandler := projectstatushandler.NewProjectStatusHandler(db)
		r.Post("/", projectStatusHandler.Add)
		r.Post("/search", projectStatusHandler.Search)
		r.Get("/{id}", projectStatusHandler.Get)
		r.Put("/{id}", projectStatusHandler.Update)
		r.Delete("/{id}", projectStatusHandler.Delete)
	})

	apiRouter.Route("/project-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectTypeHandler := projectpriorityhandler.NewProjectPriorityHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Post("/search", projectTypeHandler.Search)
		r.Get("/{id}", projectTypeHandler.Get)
		r.Put("/{id}", projectTypeHandler.Update)
		r.Delete("/{id}", projectTypeHandler.Delete)
	})

	apiRouter.Route("/task-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		taskStatusHandler := taskstatushandler.NewTaskStatusHandler(db)
		r.Post("/", taskStatusHandler.Add)
		r.Post("/search", taskStatusHandler.Search)
		r.Get("/{id}", taskStatusHandler.Get)
		r.Put("/{id}", taskStatusHandler.Update)
		r.Delete("/{id}", taskStatusHandler.Delete)
	})

	apiRouter.Route("/task-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		taskPriorityHandler := taskpriorityhandler.NewTaskPriorityHandler(db)
		r.Post("/", taskPriorityHandler.Add)
		r.Post("/search", taskPriorityHandler.Search)
		r.Get("/{id}", taskPriorityHandler.Get)
		r.Put("/{id}", taskPriorityHandler.Update)
		r.Delete("/{id}", taskPriorityHandler.Delete)
	})

	apiRouter.Route("/projects", func(r chi.Router) {
		projectHandler := projecthandler.NewProjectHandler(db)
		r.Post("/", projectHandler.AddProject)
		r.Get("/", projectHandler.SearchProjects)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", projectHandler.GetProject)
			r.Put("/", projectHandler.UpdateProject)
			r.Delete("/", projectHandler.DeleteProject)
		})
	})

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
