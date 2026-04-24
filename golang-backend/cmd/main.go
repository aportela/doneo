package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aportela/gotask/internal/cli"
	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/router"
	"github.com/aportela/gotask/internal/seed"
)

func main() {
	log.Println("starting GOTask v0.1alpha...")

	db, err := database.Open(true)
	if err != nil {
		log.Fatal(err)
	} else {
		defer db.Close()
		params, err := cli.HandleFlags()
		if err != nil {
			log.Fatal(err)
		}

		if params.InsertBulkData {
			seed.CreateDemoData(db)
		}

		r := router.NewRouter(db)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			if err := http.ListenAndServe(":3000", r); err != nil {
				log.Fatal("Error en el servidor HTTP:", err)
			}
		}()

		sigReceived := <-sigChan
		log.Printf("%v SIG received... closing app", sigReceived)
	}
}
