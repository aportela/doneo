package repositories

import (
	"github.com/aportela/doneo/internal/domain"
)

func TimestampFilterToDTO(filter *domain.TimestampFilter) *TimestampFilter {
	if filter == nil {
		return nil
	}

	return &TimestampFilter{
		From:   filter.From,
		To:     filter.To,
		Filled: filter.Filled,
		Empty:  filter.Empty,
	}
}
