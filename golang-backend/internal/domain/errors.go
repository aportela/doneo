package domain

import (
	"fmt"
)

type DatabaseError struct {
	Message string
}

func (e *DatabaseError) Error() string {
	return e.Message
}

func NewDatabaseError(msg string) error {
	return &DatabaseError{Message: msg}
}

var (
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
	ErrNotFound           = fmt.Errorf("entity not found")
	ErrDeleted            = fmt.Errorf("entity deleted")
)
