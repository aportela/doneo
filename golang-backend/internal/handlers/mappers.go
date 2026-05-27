package handlers

import (
	"errors"
	"net/http"

	"github.com/aportela/doneo/internal/domain"
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

	return http.StatusInternalServerError, "internal error", nil
}
