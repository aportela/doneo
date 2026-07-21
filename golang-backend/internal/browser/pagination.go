package browser

import (
	"math"
)

const (
	DefaultResultsPage = 20
	MaxResultsPage     = 500
)

type PagerQuery struct {
	Enabled     bool
	CurrentPage int
	ResultsPage int
}

func (p PagerQuery) getCurrentPage() int {
	if p.CurrentPage <= 0 {
		return 1
	}
	return p.CurrentPage
}

func (p PagerQuery) getResultsPage() int {
	if p.ResultsPage < 0 {
		return DefaultResultsPage
	}

	if p.ResultsPage > MaxResultsPage {
		return MaxResultsPage
	}

	return p.ResultsPage
}

func (p PagerQuery) Limit() int {
	return p.getResultsPage()
}

func (p PagerQuery) Offset() int {
	return (p.getCurrentPage() - 1) * p.getResultsPage()
}

type PagerResult struct {
	Emabled      bool
	CurrentPage  int
	ResultsPage  int
	TotalResults int
	TotalPages   int
}

func NewPagerResult(pagerQuery PagerQuery, totalResults int) PagerResult {
	resultsPage := pagerQuery.getResultsPage()
	currentPage := pagerQuery.getCurrentPage()

	totalPages := 0
	if resultsPage > 0 && totalResults > 0 {
		totalPages = int(math.Ceil(float64(totalResults) / float64(resultsPage)))
	}

	return PagerResult{
		Emabled:      pagerQuery.Enabled,
		CurrentPage:  currentPage,
		ResultsPage:  resultsPage,
		TotalResults: totalResults,
		TotalPages:   totalPages,
	}
}
