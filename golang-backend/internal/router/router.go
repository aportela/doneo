package router

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/doneo/internal/config"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/handlers/authhandler"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
	"github.com/aportela/doneo/internal/handlers/projectpermissionhandler"
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

	baseRouter.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "route not found",
		})
	})

	apiRouter := chi.NewRouter()

	apiRouter.Route("/auth", func(r chi.Router) {
		userHandler := authhandler.NewAuthHandler(db, cfg.Auth.SecretKey, cfg.Auth.AccessTokenExpirationHours, cfg.Auth.RefreshTokenExpirationDays)
		r.Post("/signin", userHandler.SignIn)
		r.Post("/signout", userHandler.SignOut)
		r.Post("/renew-access-token", userHandler.RenewAccessToken)
	})

	uuidPattern := "[0-9a-fA-F-]{36}"

	apiRouter.Route("/avatars", func(r chi.Router) {
		r.Get("/{size:[0-9]+}/user/{id:"+uuidPattern+"}", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(
				w,
				r,
				"https://i.pravatar.cc/"+chi.URLParam(r, "size")+"?u="+chi.URLParam(r, "id"),
				http.StatusTemporaryRedirect,
			)
		})
	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		userHandler := userhandler.NewUserHandler(db)
		r.Post("/", userHandler.Add)
		r.Post("/search", userHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", userHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", userHandler.Update)
		r.Patch("/{id:"+uuidPattern+"}", userHandler.Patch)
		r.Delete("/{id:"+uuidPattern+"}", userHandler.Delete)
	})

	apiRouter.Route("/roles", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		roleHandler := rolehandler.NewRoleHandler(db)
		r.Post("/", roleHandler.Add)
		r.Post("/search", roleHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", roleHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", roleHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", roleHandler.Delete)
	})

	apiRouter.Route("/project-types", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectTypeHandler := projecttypehandler.NewProjectTypeHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Post("/search", projectTypeHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", projectTypeHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", projectTypeHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", projectTypeHandler.Delete)
	})

	apiRouter.Route("/project-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectStatusHandler := projectstatushandler.NewProjectStatusHandler(db)
		r.Post("/", projectStatusHandler.Add)
		r.Post("/search", projectStatusHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", projectStatusHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", projectStatusHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", projectStatusHandler.Delete)
	})

	apiRouter.Route("/project-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		projectTypeHandler := projectpriorityhandler.NewProjectPriorityHandler(db)
		r.Post("/", projectTypeHandler.Add)
		r.Post("/search", projectTypeHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", projectTypeHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", projectTypeHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", projectTypeHandler.Delete)
	})

	apiRouter.Route("/task-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		taskStatusHandler := taskstatushandler.NewTaskStatusHandler(db)
		r.Post("/", taskStatusHandler.Add)
		r.Post("/search", taskStatusHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", taskStatusHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", taskStatusHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", taskStatusHandler.Delete)
	})

	apiRouter.Route("/task-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		taskPriorityHandler := taskpriorityhandler.NewTaskPriorityHandler(db)
		r.Post("/", taskPriorityHandler.Add)
		r.Post("/search", taskPriorityHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", taskPriorityHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", taskPriorityHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", taskPriorityHandler.Delete)
	})

	apiRouter.Route("/projects", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		projectHandler := projecthandler.NewProjectHandler(db)
		projectPermissionHandler := projectpermissionhandler.NewProjectPermissionHandler(db)
		r.Post("/", projectHandler.Add)
		r.Post("/search", projectHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", projectHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", projectHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", projectHandler.Delete)
		r.Get("/{id:"+uuidPattern+"}/permissions", projectPermissionHandler.Search)
		r.Delete("/{id:"+uuidPattern+"}/permissions/{permission_id:"+uuidPattern+"}", projectPermissionHandler.Delete)
	})

	// TODO: 404 route ?
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
