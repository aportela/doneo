package database

import (
	"errors"

	"github.com/aportela/doneo/internal/configuration"
	sqliteWrapper "github.com/aportela/doneo/internal/database/sqlite"
)

// TODO: DSN
func Open(cfg configuration.DatabaseConfiguration) (Database, error) {
	switch cfg.Type {

	case "sqlite":
		databaseHandler := &sqliteWrapper.Handler{}

		if err := databaseHandler.Open(cfg); err != nil {
			return nil, err
		}

		return databaseHandler, nil

	default:
		return nil, errors.New("unknown database type")
	}
}
