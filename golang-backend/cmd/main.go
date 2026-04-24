package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aportela/gotask/internal/cli"
	"github.com/aportela/gotask/internal/configuration"
	"github.com/aportela/gotask/internal/data"
	"github.com/aportela/gotask/internal/database"
	"github.com/aportela/gotask/internal/router"
	"github.com/aportela/gotask/internal/seed"
)

func main() {
	log.Println("starting GOTask v0.1alpha...")

	err := data.CreateDataPathIfRequired()
	if err != nil {
		log.Fatal("Error checking/creating data path:", err)
	}

	configuration, err := configuration.Open()

	if err != nil {
		log.Fatal("Error opening configuration:", err)
	}

	databaseHandler, err := database.Open(configuration.Database)
	if err != nil {
		log.Fatal(err)
	} else {
		defer databaseHandler.Close()
		createSchema := true
		if createSchema {
			log.Println("Creating database schema...")
			err = databaseHandler.CreateSchema()

			if err != nil {
				log.Fatal("Error creating database schema:", err)
			}
		}

		params, err := cli.HandleFlags()
		if err != nil {
			log.Fatal(err)
		}

		if params.InsertBulkData {
			seed.CreateDemoData(databaseHandler)
		}

		r := router.NewRouter(databaseHandler)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			log.Println("Starting web server at port " + strconv.Itoa(configuration.Server.Port))
			if err := http.ListenAndServe(":"+strconv.Itoa(configuration.Server.Port), r); err != nil {
				log.Fatal("Error", err)
			}
		}()

		sigReceived := <-sigChan
		log.Printf("%v SIG received... closing app", sigReceived)
	}
}
