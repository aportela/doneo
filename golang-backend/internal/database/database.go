package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {

	databaseDir := filepath.Join(".", "data")
	err := os.MkdirAll(databaseDir, os.ModePerm)
	databasePath := filepath.Join(databaseDir, "gotask.sqlite3")
	db, err := sql.Open("sqlite", "file:"+databasePath)
	if err != nil {
		log.Fatal(err)
	}

	schema := `PRAGMA journal_mode = WAL;`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	schema = `PRAGMA foreign_keys = ON;`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	schema = `
	CREATE TABLE IF NOT EXISTS PROJECT (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	schema = `
	CREATE TABLE IF NOT EXISTS TASK (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		FOREIGN KEY (project_id) REFERENCES PROJECT(id)
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	return db
}
