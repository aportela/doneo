package configuration

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/aportela/gotask/internal/data"
	"github.com/spf13/viper"
)

const configurationFilename = "configuration.yaml"

// configuration default values
const databaseType = "sqlite"
const sqliteDatabaseFilename = "gotask.sqlite3"
const httpServerPort = 8086

func initViper() {
	configFile := filepath.Join(data.GetDataPath(), configurationFilename)
	viper.SetConfigFile(configFile)
}

func createDefaultConfiguration() error {
	viper.Set("database.type", databaseType)
	viper.Set("database.path", filepath.Join(data.GetDataPath(), sqliteDatabaseFilename))
	viper.Set("server.port", httpServerPort)

	return viper.WriteConfigAs(filepath.Join(data.GetDataPath(), configurationFilename))
}

func Open() (*Configuration, error) {
	initViper()

	err := viper.ReadInConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Println("Configuration file not found, creating new with default values")
			viper.SetDefault("database.path", "test")
			err := createDefaultConfiguration()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	var cfg *Configuration
	err = viper.Unmarshal(&cfg)
	return cfg, err
}
