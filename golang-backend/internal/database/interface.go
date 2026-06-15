package database

import (
	"context"
	"database/sql"

	"github.com/aportela/doneo/internal/config"
)

type DatabaseExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Database interface {
	Open(cfg config.DatabaseConfiguration) error
	Begin() (*sql.Tx, error)
	CheckSchema() error
	Close() error
	DatabaseExecutor
}
