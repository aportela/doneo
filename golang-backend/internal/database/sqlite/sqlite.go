package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/aportela/doneo/internal/config"

	_ "modernc.org/sqlite"
)

type Handler struct {
	database       *sql.DB
	currentVersion uint16
}

func configure(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
		"PRAGMA busy_timeout = 1000;",
	}
	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("sqlite pragma error (%s): %w", p, err)
		}
	}

	return nil
}

func (handler *Handler) Open(databaseConfiguration config.DatabaseConfiguration) error {
	database, err := sql.Open("sqlite", "file:"+databaseConfiguration.Path)
	if err != nil {
		return fmt.Errorf("sqlite open error: %w", err)
	}

	if err := configure(database); err != nil {
		database.Close()
		return fmt.Errorf("sqlite configure error: %w", err)
	}

	if err := database.Ping(); err != nil {
		return fmt.Errorf("sqlite ping error: %w", err)
	}

	database.SetMaxOpenConns(1)
	handler.database = database
	return nil
}

func (handler *Handler) Begin() (*sql.Tx, error) {
	return handler.database.Begin()
}

func (handler *Handler) CheckSchema() error {
	ctx := context.Background()
	var currentVersion uint16

	err := handler.database.
		QueryRowContext(ctx, "PRAGMA user_version").
		Scan(&currentVersion)

	if err != nil {
		return fmt.Errorf("error getting current version: %w", err)
	}

	log.Printf("Current schema version: %d", currentVersion)

	tx, err := handler.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}
	defer tx.Rollback()

	for _, m := range schemaQueries {
		if m.Version <= currentVersion {
			continue
		}

		log.Printf("Applying schema version %d", m.Version)

		for _, q := range m.Queries {
			if _, err := tx.ExecContext(ctx, q); err != nil {
				return fmt.Errorf(
					"migration %d failed: %w",
					m.Version,
					err,
				)
			}
		}

		if _, err := tx.ExecContext(
			ctx,
			fmt.Sprintf("PRAGMA user_version=%d", m.Version),
		); err != nil {
			return fmt.Errorf(
				"setting schema version %d error: %w",
				m.Version,
				err,
			)
		}

		currentVersion = m.Version
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}

func (handler *Handler) Close() error {
	if handler.database == nil {
		return nil
	}
	return handler.database.Close()
}

func (handler *Handler) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return handler.database.ExecContext(ctx, query, args...)
}

func (handler *Handler) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return handler.database.QueryContext(ctx, query, args...)
}

func (handler *Handler) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return handler.database.QueryRowContext(ctx, query, args...)
}
