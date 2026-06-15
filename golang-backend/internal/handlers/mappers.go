package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
	"modernc.org/sqlite"
)

func ToEmptyResponse() EmptyResponse {
	return EmptyResponse{}
}

func mapError(err error) (int, string, any) {

	if errors.Is(err, domain.NotFoundError) {
		return http.StatusNotFound, "resource not found", nil
	} else if errors.Is(err, domain.DeletedError) {
		return http.StatusGone, "resource has been deleted", nil
	} else if errors.Is(err, domain.InvalidCredentialsError) {
		return http.StatusUnauthorized, "invalid credentials", nil
	} else if errors.Is(err, domain.AuthorizationError) {
		return http.StatusForbidden, "access denied", nil
	}

	var alreadyExistsError *domain.AlreadyExistsError
	if errors.As(err, &alreadyExistsError) {
		return http.StatusConflict, "resource already exists", map[string]string{
			"field": alreadyExistsError.Field,
		}
	}

	var validationError *domain.ValidationError
	if errors.As(err, &validationError) {
		return http.StatusBadRequest, "bad request", map[string]string{
			"field": validationError.Field,
		}
	}

	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		// print SQL to console
		fmt.Println(err.Error())
	}

	return http.StatusInternalServerError, "internal error", nil
}
