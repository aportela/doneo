package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/aportela/doneo/internal/cli"
	"github.com/aportela/doneo/internal/configuration"
	"github.com/aportela/doneo/internal/data"
	"github.com/aportela/doneo/internal/database"
	"github.com/aportela/doneo/internal/router"
	"github.com/aportela/doneo/internal/scripts"
	demodatascripts "github.com/aportela/doneo/internal/scripts/demo_data"
)

func main() {
	log.Println("starting Doneo v0.1alpha...")

	err := data.CreateDataPathIfRequired()
	if err != nil {
		log.Fatal("Error checking/creating data path:", err)
	}

	configuration, err := configuration.Open()

	if err != nil {
		log.Fatal("Error opening configuration:", err)
	}

	_, err = os.Stat(configuration.Database.Path)
	createSchema := false
	if err != nil {
		createSchema = true
	}

	databaseHandler, err := database.Open(configuration.Database)
	if err != nil {
		log.Fatal(err)
	} else {
		defer databaseHandler.Close()
		if createSchema {
			log.Println("Creating database schema...")
			err = databaseHandler.CreateSchema()

			if err != nil {
				log.Fatal("Error creating database schema:", err)
			}
			scripts.CreateDefaultData(databaseHandler)
		}

		params, err := cli.HandleFlags()
		if err != nil {
			log.Fatal(err)
		}

		if params.InsertBulkData {
			log.Println("Inserting bulk/demo data...")
			demodatascripts.CreateDemoData(databaseHandler)
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
