package configuration

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configurationFilename = "configuration"
const configurationType = "yaml"

func getDataPath() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting data directory:", err)
	}
	return filepath.Join(pwd, "data")
}

func createDataPathIfRequired() error {
	if err := os.MkdirAll(getDataPath(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func initViper() {
	viper.AddConfigPath(getDataPath())
	viper.SetConfigName(configurationFilename)
	viper.SetConfigType(configurationType)
}

func Init() {
	// TODO: return error (replace log.Fatal)
	err := createDataPathIfRequired()
	if err != nil {
		log.Fatal("Error checking data path:", err)
	}
	initViper()

	err = viper.ReadInConfig()
	if err != nil {
		var notFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &notFoundErr) {
			log.Println("Configuration file not found, creating new with default values")
			viper.SetDefault("database.path", "test")
			err := writeDefaultConfig()
			if err != nil {
				log.Fatal("Error creating configuration file:", err)
			}
		} else {
			log.Fatal("Error reading configuration file:", err)
		}
	}
	dbPath := viper.GetString("database.path")
	fmt.Println("Database path:", dbPath)
}

func writeDefaultConfig() error {
	viper.Set("database.path", filepath.Join(getDataPath(), "gotask.sqlite3"))

	err := viper.SafeWriteConfig()
	if err != nil {
		return err
	}

	return nil
}
