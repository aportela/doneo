package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func Configure(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return err
		}
	}

	return nil
}

func Open(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+databasePath)
	if err != nil {
		return nil, err
	}

	err = Configure(db)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *sql.DB) error {
	return db.Close()
}
