package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound         = fmt.Errorf("entity not found")
	ErrSQLDatabaseError = errors.New("SQL database error: ")
)

func WrapSQLError(err error) error {
	return fmt.Errorf("%w: %v", ErrSQLDatabaseError, err)
}
