package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aportela/gotask/internal/configuration"

	_ "modernc.org/sqlite"
)

type Handler struct {
	database *sql.DB
}

func configure(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA foreign_keys = ON;",
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			return fmt.Errorf("sqlite pragma error (%s): %w", p, err)
		}
	}

	return nil
}

func (handler *Handler) Open(databaseConfiguration configuration.DatabaseConfiguration) error {
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

	handler.database = database
	return nil
}

func (handler *Handler) CreateSchema() error {

	tx, err := handler.database.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	for _, query := range installSchemaQueries {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return fmt.Errorf("sqlite schema error: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
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
