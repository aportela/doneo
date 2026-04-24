package configuration

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const configurationType = "yaml"
const configurationFilename = "configuration." + configurationType

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
	viper.SetConfigType(configurationType)
	viper.SetConfigFile(filepath.Join(getDataPath(), configurationFilename))

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
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("Error decoding configuration file:", err)
	}
}

func writeDefaultConfig() error {
	viper.Set("database.type", "sqlite")
	viper.Set("database.path", filepath.Join(getDataPath(), "gotask.sqlite3"))
	viper.Set("server.port", 3000)

	err := viper.SafeWriteConfig()
	if err != nil {
		return err
	}

	return nil
}
