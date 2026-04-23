package main

import (
	"log"
	"net/http"

	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/router"
)

func main() {
	log.Println("starting GOTask v0.1alpha...")

	db, err := database.Open(true)
	if err != nil {
		log.Fatal(err)
	} else {
		db.Close()

		r := router.NewRouter()

		log.Println("Listening over http://localhost:3000/static")

		http.ListenAndServe(":3000", r)
	}
}
