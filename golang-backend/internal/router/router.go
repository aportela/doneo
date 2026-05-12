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

	apiRouter.Route("/project-types", func(r chi.Router) {
		projectTypeHandler := projecttypehandler.NewProjectTypeHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Get("/", projectTypeHandler.Search)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", projectTypeHandler.Get)
			r.Put("/", projectTypeHandler.Update)
			r.Delete("/", projectTypeHandler.Delete)
		})
	})

	apiRouter.Route("/project-statuses", func(r chi.Router) {
		projectStatusHandler := projectstatushandler.NewProjectStatusHandler(db)
		r.Post("/", projectStatusHandler.Add)
		r.Get("/", projectStatusHandler.Search)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", projectStatusHandler.Get)
			r.Put("/", projectStatusHandler.Update)
			r.Delete("/", projectStatusHandler.Delete)
		})
	})

	apiRouter.Route("/project-priorities", func(r chi.Router) {
		projectTypeHandler := projectpriorityhandler.NewProjectPriorityHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Get("/", projectTypeHandler.Search)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", projectTypeHandler.Get)
			r.Put("/", projectTypeHandler.Update)
			r.Delete("/", projectTypeHandler.Delete)
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
