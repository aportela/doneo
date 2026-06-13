package domain

import (
	"errors"
)

var AuthorizationError = errors.New("access denied")
var InvalidCredentialsError = errors.New("invalid credentials")
var NotFoundError = errors.New("entity not found")
var DeletedError = errors.New("entity deleted")
var UncaughtDatabaseError = errors.New("database error")

type AlreadyExistsError struct {
	Field string
}

func (e *AlreadyExistsError) Error() string {
	return "field already exists: " + e.Field
}

type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "field not valid: " + e.Field
}
