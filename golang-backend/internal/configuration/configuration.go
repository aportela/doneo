package configuration

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/aportela/doneo/internal/data"
	"github.com/spf13/viper"
)

const configurationFilename = "configuration.yaml"

// configuration default values
const databaseType = "sqlite"
const sqliteDatabaseFilename = "doneo.sqlite3"
const httpServerPort = 8086

const accessTokenExpirationDays = 1
const refreshTokenExpirationDays = 365

func initViper() {
	configFile := filepath.Join(data.GetDataPath(), configurationFilename)
	viper.SetConfigFile(configFile)
}

func generateSecretKey(length int) (string, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func createDefaultConfiguration() error {
	viper.Set("database.type", databaseType)
	viper.Set("database.path", filepath.Join(data.GetDataPath(), sqliteDatabaseFilename))
	viper.Set("server.port", httpServerPort)

	viper.Set("auth.access_token_expiration_days", accessTokenExpirationDays)
	viper.Set("auth.refresh_token_expiration_days", refreshTokenExpirationDays)
	secretKey, err := generateSecretKey(128)
	if err != nil {
		return err
	}
	viper.Set("auth.secret_key", secretKey)

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
