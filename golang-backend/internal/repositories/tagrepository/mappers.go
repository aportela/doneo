package tagrepository

import (
	"errors"
	"strings"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}
	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(tag) ") {
			return &domain.ValidationError{Field: "tag"}
		}
	}
	return err
}
