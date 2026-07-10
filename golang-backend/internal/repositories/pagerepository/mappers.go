package pagerepository

import (
	"errors"
	"strings"
	"time"

	"github.com/aportela/doneo/internal/domain"
	"github.com/aportela/doneo/internal/utils"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func toDTO(page domain.Page) pageDTO {
	return pageDTO{
		ID:          page.ID,
		CreatorID:   page.CreatedBy.ID,
		CreatorName: page.CreatedBy.Name,
		CreatedAt:   page.CreatedAt.UnixMilli(),
		UpdatedAt:   utils.TimePtrToSQLNullInt64(page.UpdatedAt),
		Title:       page.Title,
		Body:        page.Body,
	}
}

func toDomain(page pageDTO) domain.Page {
	return domain.Page{
		ID: page.ID,
		CreatedBy: domain.UserBase{
			ID:   page.CreatorID,
			Name: page.CreatorName,
		},
		CreatedAt: time.UnixMilli(page.CreatedAt),
		UpdatedAt: utils.SQLNullInt64ToTimePtr(page.UpdatedAt),
		Title:     page.Title,
		Body:      page.Body,
	}
}

func toDomainArray(pages []pageDTO) []domain.Page {
	results := make([]domain.Page, 0, len(pages))
	for _, page := range pages {
		results = append(results, toDomain(page))
	}
	return results
}

func mapSQLiteError(err error) error {
	var sqlErr *sqlite.Error
	if !errors.As(err, &sqlErr) {
		return err
	}

	switch sqlErr.Code() {
	case sqlite3.SQLITE_CONSTRAINT_CHECK:
		if strings.Contains(sqlErr.Error(), "length(title)") {
			return &domain.ValidationError{Field: "title"}
		}
	}

	return err
}
