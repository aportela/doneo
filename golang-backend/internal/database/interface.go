package database

import (
	"context"
	"database/sql"

	"github.com/aportela/gotask/internal/configuration"
)

type Database interface {
	Open(cfg configuration.DatabaseConfiguration) error
	CreateSchema() error
	Close() error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
