package router

import (
	"encoding/json"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/doneo/internal/app"
	"github.com/aportela/doneo/internal/middlewares"

	"github.com/aportela/doneo/internal/ui"
)

func NewRouter(app *app.App) http.Handler {
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
		r.Post("/signin", app.AuthHandler.SignIn)
		r.Post("/signout", app.AuthHandler.SignOut)
		r.Post("/renew-access-token", app.AuthHandler.RenewAccessToken)
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
		r.Use(middlewares.RequireJWTCookieAuthentication(app.Cfg.Auth.SecretKey))
		r.Get("/project/{id:"+uuidPattern+"}/attachment/{attachment_id:"+uuidPattern+"}", app.AttachmentHandler.DownloadProjectAttachment)

	})

	apiRouter.Route("/entities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Get("/users", app.UserHandler.SearchBase)
		r.Get("/roles", app.RoleHandler.SearchBase)
	})

	apiRouter.Route("/users", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.UserHandler.Add)
		r.Post("/search", app.UserHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.UserHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.UserHandler.Update)
		r.Patch("/{id:"+uuidPattern+"}", app.UserHandler.Patch)
		r.Delete("/{id:"+uuidPattern+"}", app.UserHandler.Delete)
	})

	apiRouter.Route("/roles", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.RoleHandler.Add)
		r.Post("/search", app.RoleHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.RoleHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.RoleHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.RoleHandler.Delete)
	})

	apiRouter.Route("/project-types", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.ProjectTypeHandler.Add)
		r.Post("/search", app.ProjectTypeHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.ProjectTypeHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.ProjectTypeHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.ProjectTypeHandler.Delete)
	})

	apiRouter.Route("/project-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.ProjectStatusHandler.Add)
		r.Post("/search", app.ProjectStatusHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.ProjectStatusHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.ProjectStatusHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.ProjectStatusHandler.Delete)
	})

	apiRouter.Route("/project-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.ProjectPriorityHandler.Add)
		r.Post("/search", app.ProjectPriorityHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.ProjectPriorityHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.ProjectPriorityHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.ProjectPriorityHandler.Delete)
	})

	apiRouter.Route("/task-statuses", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.TaskStatusHandler.Add)
		r.Post("/search", app.TaskStatusHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.TaskStatusHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.TaskStatusHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.TaskStatusHandler.Delete)
	})

	apiRouter.Route("/task-priorities", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Use(middlewares.RequireSuperUser)
		r.Post("/", app.TaskPriorityHandler.Add)
		r.Post("/search", app.TaskPriorityHandler.Search)
		r.Get("/{id:"+uuidPattern+"}", app.TaskPriorityHandler.Get)
		r.Put("/{id:"+uuidPattern+"}", app.TaskPriorityHandler.Update)
		r.Delete("/{id:"+uuidPattern+"}", app.TaskPriorityHandler.Delete)
	})

	apiRouter.Route("/timers", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Post("/", app.TimerHandler.Start)
		r.Put("/{id:"+uuidPattern+"}", app.TimerHandler.Stop)
		r.Delete("/{id:"+uuidPattern+"}", app.TimerHandler.Delete)
		r.Delete("/", app.TimerHandler.Clear)
		r.Get("/", app.TimerHandler.Search)
	})

	apiRouter.Route("/projects", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Post("/", app.ProjectHandler.Add)
		r.Post("/search", app.ProjectHandler.Search)
		r.Get("/{project_id:"+uuidPattern+"}", app.ProjectHandler.Get)
		r.Put("/{project_id:"+uuidPattern+"}", app.ProjectHandler.Update)
		r.Delete("/{project_id:"+uuidPattern+"}", app.ProjectHandler.Delete)

		r.Get("/{project_id:"+uuidPattern+"}/permissions", app.ProjectPermissionHandler.Search)
		r.Post("/{project_id:"+uuidPattern+"}/permissions", app.ProjectPermissionHandler.Add)
		r.Delete("/{project_id:"+uuidPattern+"}/permissions/{permission_id:"+uuidPattern+"}", app.ProjectPermissionHandler.Delete)

		r.Get("/{project_id:"+uuidPattern+"}/notes", app.NoteHandler.GetProjectNotes)
		r.Post("/{project_id:"+uuidPattern+"}/notes", app.NoteHandler.AddProjectNote)
		r.Put("/{project_id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", app.NoteHandler.UpdateProjectNote)
		r.Delete("/{project_id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", app.NoteHandler.DeleteProjectNote)

		r.Get("/{project_id:"+uuidPattern+"}/attachments", app.AttachmentHandler.GetProjectAttachments)
		r.Post("/{project_id:"+uuidPattern+"}/attachments", app.AttachmentHandler.AddProjectAttachment)
		r.Delete("/{project_id:"+uuidPattern+"}/attachments/{attachment_id:"+uuidPattern+"}", app.AttachmentHandler.DeleteProjectAttachment)

		r.Get("/{project_id:"+uuidPattern+"}/history_operations", app.HistoryOperationHandler.SearchProjectHistoryOperations)

		r.Post("/{project_id:"+uuidPattern+"}/tasks/search", app.TaskHandler.Search)
		r.Post("/{project_id:"+uuidPattern+"}/tasks", app.TaskHandler.Add)
		r.Get("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}", app.TaskHandler.Get)
		r.Put("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}", app.TaskHandler.Update)
		r.Delete("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}", app.TaskHandler.Delete)

		r.Get("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/notes", app.NoteHandler.GetTaskNotes)
		r.Post("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/notes", app.NoteHandler.AddTaskNote)
		r.Put("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", app.NoteHandler.UpdateTaskNote)
		r.Delete("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/notes/{note_id:"+uuidPattern+"}", app.NoteHandler.DeleteTaskNote)

		r.Get("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/attachments", app.AttachmentHandler.GetTaskAttachments)
		r.Post("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/attachments", app.AttachmentHandler.AddTaskAttachment)
		r.Delete("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/attachments/{attachment_id:"+uuidPattern+"}", app.AttachmentHandler.DeleteTaskAttachment)

		r.Get("/{project_id:"+uuidPattern+"}/tasks/{task_id:"+uuidPattern+"}/history_operations", app.HistoryOperationHandler.SearchTaskHistoryOperations)
	})

	apiRouter.Route("/tasks", func(r chi.Router) {
		r.Use(middlewares.RequireJWTAuthentication(app.Cfg.Auth.SecretKey))
		r.Post("/search", app.TaskHandler.Search)
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
