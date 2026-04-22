package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/aportela/gotask/internal/fileserver"
	"github.com/aportela/gotask/internal/handlers"
)

func NewRouter() http.Handler {
	baseRouter := chi.NewRouter()

	baseRouter.Use(middleware.Logger)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/hello", handlers.DefaultHandler)
	baseRouter.Mount("/api", apiRouter)

	workDir, _ := os.Getwd()
	publicWebPath := http.Dir(filepath.Join(workDir, "public"))
	fileserver.FileServer(baseRouter, "/", publicWebPath)

	return baseRouter
}
