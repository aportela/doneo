package main

import (
	"log"
	"net/http"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/router"
)

func main() {

	log.Println("starting GOTask v0.1alpha...")

	db := database.InitDB()

	r := router.NewRouter()

	http.ListenAndServe(":3000", r)

}
