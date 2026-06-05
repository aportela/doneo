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
	"github.com/aportela/doneo/internal/handlers/attachmenthandler"
	"github.com/aportela/doneo/internal/handlers/authhandler"
	"github.com/aportela/doneo/internal/handlers/notehandler"
	"github.com/aportela/doneo/internal/handlers/projecthandler"
	"github.com/aportela/doneo/internal/handlers/projecthistoryhandler"
	"github.com/aportela/doneo/internal/handlers/projectpermissionhandler"
	"github.com/aportela/doneo/internal/handlers/projectpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/projectstatushandler"
	"github.com/aportela/doneo/internal/handlers/projecttypehandler"
	"github.com/aportela/doneo/internal/handlers/rolehandler"
	"github.com/aportela/doneo/internal/handlers/taskpriorityhandler"
	"github.com/aportela/doneo/internal/handlers/taskstatushandler"
	"github.com/aportela/doneo/internal/handlers/userhandler"
	"github.com/aportela/doneo/internal/middlewares"
	"github.com/aportela/doneo/internal/repositories/projecthistoryrespository"
	"github.com/aportela/doneo/internal/repositories/projectpermissionrepository"
	"github.com/aportela/doneo/internal/repositories/projectpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/projectrepository"
	"github.com/aportela/doneo/internal/repositories/projectstatusrepository"
	"github.com/aportela/doneo/internal/repositories/projecttyperepository"
	"github.com/aportela/doneo/internal/repositories/rolerepository"
	"github.com/aportela/doneo/internal/repositories/taskpriorityrepository"
	"github.com/aportela/doneo/internal/repositories/taskstatusrepository"
	"github.com/aportela/doneo/internal/repositories/userrepository"
	"github.com/aportela/doneo/internal/services/authservice"
	"github.com/aportela/doneo/internal/services/projecthistoryservice"
	"github.com/aportela/doneo/internal/services/projectpermissionservice"
	"github.com/aportela/doneo/internal/services/projectpriorityservice"
	"github.com/aportela/doneo/internal/services/projectservice"
	"github.com/aportela/doneo/internal/services/projectstatusservice"
	"github.com/aportela/doneo/internal/services/projecttypeservice"
	"github.com/aportela/doneo/internal/services/roleservice"
	"github.com/aportela/doneo/internal/services/taskpriorityservice"
	"github.com/aportela/doneo/internal/services/taskstatusservice"
	"github.com/aportela/doneo/internal/services/userservice"

	"github.com/aportela/doneo/internal/ui"
)

func NewRouter(database database.Database, cfg config.Configuration) http.Handler {
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
		handler := authhandler.NewHandler(authservice.NewService(database, userrepository.NewRepository(database)), cfg.Auth.SecretKey, cfg.Auth.AccessTokenExpirationHours, cfg.Auth.RefreshTokenExpirationDays)
		r.Post("/signin", handler.SignIn)
		r.Post("/signout", handler.SignOut)
		r.Post("/renew-access-token", handler.RenewAccessToken)
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

	apiRouter.Route("/attachments", func(r chi.Router) {
		r.Use(middlewares.RequireJWTCookieAuthentication(cfg.Auth.SecretKey))
		// TODO: remove project, check attachment permissions on get
		handler := attachmenthandler.NewHandler(database, cfg.Storage.AttachmentsPath)
		r.Get("/project/{id:"+uuidPattern+"}/attachment/{attachment_id:"+uuidPattern+"}", handler.DownloadProjectAttachment)

	})

	apiRouter.Route("/entities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		userHandler := userhandler.NewHandler(userservice.NewService(database, userrepository.NewRepository(database)))
		roleHandler := rolehandler.NewHandler(roleservice.NewService(database, rolerepository.NewRepository(database)))
		r.Get("/users", userHandler.SearchBase)
		r.Get("/roles", roleHandler.SearchBase)
	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := userhandler.NewHandler(userservice.NewService(database, userrepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Patch("/{id:"+uuidPattern+"}", handler.Patch)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/roles", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := rolehandler.NewHandler(roleservice.NewService(database, rolerepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/project-types", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := projecttypehandler.NewHandler(projecttypeservice.NewService(database, projecttyperepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/project-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := projectstatushandler.NewHandler(projectstatusservice.NewService(database, projectstatusrepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/project-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := projectpriorityhandler.NewHandler(projectpriorityservice.NewService(database, projectpriorityrepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/task-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := taskstatushandler.NewHandler(taskstatusservice.NewService(database, taskstatusrepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/task-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		handler := taskpriorityhandler.NewHandler(taskpriorityservice.NewService(database, taskpriorityrepository.NewRepository(database)))
		r.Post("/", handler.Add)
		r.Post("/search", handler.Search)
		r.Get("/{id:"+uuidPattern+"}", handler.Get)
		r.Put("/{id:"+uuidPattern+"}", handler.Update)
		r.Delete("/{id:"+uuidPattern+"}", handler.Delete)
	})

	apiRouter.Route("/projects", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(cfg.Auth.SecretKey))
		projectHandler := projecthandler.NewHandler(projectservice.NewService(database, projectrepository.NewRepository(database)))
		projectPermissionHandler := projectpermissionhandler.NewHandler(projectpermissionservice.NewService(database, projectpermissionrepository.NewRepository(database)))
		projectNoteHandler := notehandler.NewHandler(database)
		projectAttachmentHandler := attachmenthandler.NewHandler(database, cfg.Storage.AttachmentsPath)
		projectHistoryHandler := projecthistoryhandler.NewHandler(projecthistoryservice.NewService(database, projecthistoryrespository.NewRepository(database)))
		r.Post("/", projectHandler.Add)
		r.Post("/search", projectHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", projectHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", projectHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", projectHandler.Delete)

		r.Get("/{id:"+uuidPattern+"}/permissions", projectPermissionHandler.Search)
		r.Post("/{id:"+uuidPattern+"}/permissions/", projectPermissionHandler.Add)
		r.Delete("/{id:"+uuidPattern+"}/permissions/{permission_id:"+uuidPattern+"}", projectPermissionHandler.Delete)

		r.Get("/{id:"+uuidPattern+"}/notes", projectNoteHandler.GetProjectNotes)
		r.Post("/{id:"+uuidPattern+"}/notes/", projectNoteHandler.AddProjectNote)
		r.Put("/{id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", projectNoteHandler.UpdateProjectNote)
		r.Delete("/{id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", projectNoteHandler.DeleteProjectNote)

		r.Get("/{id:"+uuidPattern+"}/attachments", projectAttachmentHandler.GetProjectAttachments)
		r.Post("/{id:"+uuidPattern+"}/attachments/", projectAttachmentHandler.AddProjectAttachment)
		r.Delete("/{id:"+uuidPattern+"}/attachments/{attachment_id:"+uuidPattern+"}", projectAttachmentHandler.DeleteProjectAttachment)

		r.Get("/{id:"+uuidPattern+"}/history_operations", projectHistoryHandler.Search)
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
