package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(path string) *sql.DB {
	db, err := sql.Open("sqlite", "file:tasks.db")
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
